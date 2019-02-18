package main

import (
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/supergiant/analyze/pkg/kube"
	"github.com/supergiant/analyze/pkg/models"
	"io/ioutil"
	"k8s.io/apimachinery/pkg/labels"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

// TODO this is copypaste from sunsetting plugin repository
type PluginInfo struct {
	// detailed plugin description
	Description string `json:"description,omitempty"`

	// unique ID of installed plugin
	// basically it is slugged URI of plugin repository name e. g. supergiant-request-limits-check
	//
	ID string `json:"id,omitempty"`

	// date/Time the plugin was installed
	// Filled by post-install job
	InstalledAt string `json:"installedAt,omitempty"`

	// name is the name of the plugin.
	Name string `json:"name,omitempty"`

	// service labels
	// Filled by post-install job
	ServiceLabels map[string]string `json:"serviceLabels,omitempty"`

	// name of k8s service which is front of plugin deployment
	// Filled by post-install job
	ServiceEndpoint string `json:"serviceEndpoint,omitempty"`

	// entry points for web components
	SettingsComponentEntryPoint string `json:"settingsComponentEntryPoint,omitempty"`
	CheckComponentEntryPoint    string `json:"checkComponentEntryPoint,omitempty"`

	// plugin status
	Status string `json:"status,omitempty"`

	// plugin version, major version shall be equal to analyze-core version
	Version string `json:"version,omitempty"`

	Revision string `json:"revision,omitempty"`
	Branch    string `json:"branch,omitempty"`
	BuildDate string `json:"buildDate,omitempty"`
	GoVersion string `json:"goVersion,omitempty"`
}

var analyzeLabelSet = labels.Set{
	"app.kubernetes.io/name": "analyze",
}

var pluginLabelSet = labels.Set{
	"app.kubernetes.io/part-of": "analyze",
	"app.kubernetes.io/component": "analyze-plugin",
}

func main() {
	command := &cobra.Command{
		Use:          "analyze-registry-job",
		Short:        "analyze-registry-job is job which registers or removes plugin from analyze registry",
		RunE:         runCommand,
		SilenceUsage: true,
	}

	command.PersistentFlags().BoolP(
		"remove",
		"r",
		false,
		"if true job will try to remove plugin from analyze registry")

	if err := command.Execute(); err != nil {
		log.Fatalf("\n%v\n", err)
	}
}

func runCommand(cmd *cobra.Command, _ []string) error {
	var rawPluginInfo []byte
	remove, err := cmd.Flags().GetBool("remove")
	if err != nil {
		return errors.Wrap(err, "unable to get config flag remove")
	}

	pluginServiceName, err := discoverPluginServiceName()
	if err != nil {
		return errors.Wrap(err, "unable to get plugin service name")
	}

	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)

	logger.Debugf("remove: %v", remove)
	logger.Debugf("pluginServiceName: %v", pluginServiceName)

	kubeClient, err := kube.NewKubeClient(logger.WithField("component", "kubeClient"))
	if err != nil {
		return errors.Wrap(err, "unable to create kube client")
	}

	pluginService, err := kubeClient.GetService(pluginServiceName, pluginLabelSet)
	if err != nil {
		return errors.Wrap(err, "failed to find analyze service")
	}

	var pluginApiPort string
	for _, port := range pluginService.Spec.Ports {
		if port.Name == "http" {
			pluginApiPort = strconv.Itoa(int(port.Port))
			break
		}
	}

	if pluginApiPort == "" {
		logger.Debugf("pluginService spec: %+v", pluginService.Spec)
		return errors.New("failed to find http port for analyze plugin")
	}


	// TODO: make it configurable
	for i := 0;i <10; i++ {
		resp, err := http.Get("http://" + pluginServiceName + ":" + pluginApiPort + "/api/v1/info")
		if err != nil || (resp != nil && resp.StatusCode != http.StatusOK) {
			logger.Debugf("unable to get plugin info: %v, statusCode: %v try in 1 sec", err, resp)
			continue
		}
		rawPluginInfo, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			logger.Debugf("unable to read plugin version from response body: %v", err)
		}
		resp.Body.Close()
		break
	}

	pi := &PluginInfo{}
	err = json.Unmarshal(rawPluginInfo, pi)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal plugin info")
	}

	logger.Debugf("plugin info: %v", string(rawPluginInfo))

	analyzeService, err := kubeClient.GetServiceByLabels(analyzeLabelSet)
	if err != nil {
		return errors.New("failed to find analyze service")
	}

	var analyzeApiPort, analyzeServiceName = "", ""

	for _, port := range analyzeService.Spec.Ports {
		if port.Name == "http" {
			analyzeApiPort = strconv.Itoa(int(port.Port))
			break
		}
	}

	if analyzeApiPort == "" {
		return errors.New("failed to find http port for analyze service")
	}

	analyzeServiceName = analyzeService.Name

	logger.Debugf("analyze service name %v, service port %v", analyzeServiceName, analyzeApiPort)

	pi.ServiceEndpoint = pluginServiceName + ":" + pluginApiPort
	pi.ServiceLabels = pluginService.Labels
	bytes, err := json.Marshal(pi)
	if err != nil {
		return errors.Wrap(err, "failed to marshal plugin info")
	}

	pluginsEndpointUri := "http://" + analyzeServiceName + ":" + analyzeApiPort + "/api/v1/plugins"
	resp, err := http.Post( pluginsEndpointUri, "application/json", strings.NewReader(string(bytes)))
	if err != nil {
		return errors.Wrap(err, "failed to register plugin")
	}

	if !(resp.StatusCode == http.StatusCreated || resp.StatusCode == http.StatusOK) {
		return errors.Errorf("failed to register plugin, status code %v", resp.StatusCode)
	}

	resp, err = http.Get(pluginsEndpointUri)
	if err != nil {
		return errors.Wrap(err, "unable to check registered plugin")
	}

	result := make([]*models.Plugin, 0, 0)
	//TODO: generate swagger client?

	defer func() {
		var err = resp.Body.Close()
		if err != nil {
			logger.Error("post request body read error")
		}
	}()

	bytes, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(err, "unable to read registered plugins")
	}

	err = json.Unmarshal(bytes, &result)
	if err != nil {
		return errors.Wrap(err, "unable to unmarshal registered plugins")
	}

	logger.Debugf("analyze plugins %v", string(bytes))
	for _, p := range result {
		logger.Debugf("analyze plugins %+v", p)
	}

	return nil
}

func discoverPluginServiceName() (string, error) {
	address, exists := os.LookupEnv("PLUGIN_SERVICE_NAME")
	if !exists {
		return "", errors.New("Environment variable PLUGIN_SERVICE_NAME is not set")
	}
	return address, nil
}