package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/labels"

	"github.com/supergiant/analyze/pkg/kube"
	"github.com/supergiant/analyze/pkg/logger"
	"github.com/supergiant/analyze/pkg/models"
)

// TODO this is copy-paste from sunsetting plugin repository
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

	Revision  string `json:"revision,omitempty"`
	Branch    string `json:"branch,omitempty"`
	BuildDate string `json:"buildDate,omitempty"`
	GoVersion string `json:"goVersion,omitempty"`
}

var analyzeLabelSet = labels.Set{
	"app.kubernetes.io/name": "analyze",
}

var pluginLabelSet = labels.Set{
	"app.kubernetes.io/part-of":   "analyze",
	"app.kubernetes.io/component": "analyze-plugin",
}

var (
	remove    = flag.Bool("remove", false, "if true job will try to remove plugin from analyze registry")
	logLevel  = flag.String("log-level", "debug", "logging level, e.g. info, warning, debug, error, fatal")
	logFormat = flag.String("log-format", "TXT", "logging format [TXT JSON]")
)

func main() {
	var rawPluginInfo []byte
	flag.Parse()

	loggerConf := logger.Config{
		Level:     *logLevel,
		Formatter: logger.Formatter(*logFormat),
	}

	if err := loggerConf.Validate(); err != nil {
		log.Fatalf("\n%v\n", err)
	}
	logger := logger.NewLogger(loggerConf).WithField("app", "analyze-registry-job")

	pluginServiceName, err := discoverPluginServiceName()
	if err != nil {
		logger.Fatalf("unable to get plugin service name, err: %v", err)
	}

	logger.Debugf("remove: %v", *remove)
	logger.Debugf("pluginServiceName: %v", pluginServiceName)

	kubeClient, err := kube.NewKubeClient(logger.WithField("component", "kubeClient"))
	if err != nil {
		logger.Fatalf("unable to create kube client, err: %v", err)
	}

	pluginService, err := kubeClient.GetService(pluginServiceName, pluginLabelSet)
	if err != nil {
		logger.Fatalf("failed to find analyze service, err: %v", err)
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
		logger.Fatalf("failed to find http port for analyze plugin")
	}

	// TODO: make it configurable
	for i := 0; i < 10; i++ {
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
		logger.Fatalf("failed to unmarshal plugin info, err: %v", err)
	}

	logger.Debugf("plugin info: %v", string(rawPluginInfo))

	analyzeService, err := kubeClient.GetServiceByLabels(analyzeLabelSet)
	if err != nil {
		logger.Fatalf("failed to find analyze service")
	}

	var analyzeApiPort, analyzeServiceName = "", ""

	for _, port := range analyzeService.Spec.Ports {
		if port.Name == "http" {
			analyzeApiPort = strconv.Itoa(int(port.Port))
			break
		}
	}

	if analyzeApiPort == "" {
		logger.Fatalf("failed to find http port for analyze service")
	}

	analyzeServiceName = analyzeService.Name

	logger.Debugf("analyze service name %v, service port %v", analyzeServiceName, analyzeApiPort)

	pi.ServiceEndpoint = pluginServiceName + ":" + pluginApiPort
	pi.ServiceLabels = pluginService.Labels
	bytes, err := json.Marshal(pi)
	if err != nil {
		logger.Fatalf("failed to marshal plugin info, err: %v", err)
	}

	pluginsEndpointUri := "http://" + analyzeServiceName + ":" + analyzeApiPort + "/api/v1/plugins"
	resp, err := http.Post(pluginsEndpointUri, "application/json", strings.NewReader(string(bytes)))
	if err != nil {
		logger.Fatalf("failed to register plugin, err: %v", err)
	}

	if !(resp.StatusCode == http.StatusCreated || resp.StatusCode == http.StatusOK) {
		logger.Fatalf("failed to register plugin, status code %v", resp.StatusCode)
	}

	resp, err = http.Get(pluginsEndpointUri)
	if err != nil {
		logger.Fatalf("unable to check registered plugin, err: %v", err)
	}

	result := make([]*models.Plugin, 0)
	//TODO: generate swagger client?

	defer func() {
		var err = resp.Body.Close()
		if err != nil {
			logger.Error("post request body read error")
		}
	}()

	bytes, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Fatalf("unable to read registered plugins, err: %v", err)
	}

	err = json.Unmarshal(bytes, &result)
	if err != nil {
		logger.Fatalf("unable to unmarshal registered plugins, err: %v", err)
	}

	logger.Debugf("analyze plugins %v", string(bytes))
	for _, p := range result {
		logger.Debugf("analyze plugins %+v", p)
	}
}

func discoverPluginServiceName() (string, error) {
	address, exists := os.LookupEnv("PLUGIN_SERVICE_NAME")
	if !exists {
		return "", errors.New("Environment variable PLUGIN_SERVICE_NAME is not set")
	}
	return address, nil
}
