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

var AnalyzeLabelSet = labels.Set{
	"app.kubernetes.io/name": "analyze",
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
	var pluginInfo = ""
	remove, err := cmd.Flags().GetBool("remove")
	if err != nil {
		return errors.Wrap(err, "unable to get config flag remove")
	}

	pluginApiAddress := discoverPluginApiAddress()

	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)

	logger.Debugf("remove: %v", remove)
	logger.Debugf("pluginApiAddress: %v", pluginApiAddress)
	// TODO: make it configurable
	for i := 0;i <10; i++ {
		resp, err := http.Get("http://" + pluginApiAddress + "/api/v1/info")
		if err != nil || (resp != nil && resp.StatusCode != http.StatusOK) {
			logger.Debugf("unable to get plugin info: %v, statusCode: %v try in 1 sec", err, resp)
			continue
		}
		bytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			logger.Debugf("unable to read plugin version from response body: %v", err)
		}
		resp.Body.Close()

		pluginInfo = string(bytes)
		break
	}

	logger.Debugf("plugin version: %v", pluginInfo)
	if pluginInfo == "" {
		return errors.New("failed to get plugin version")
	}

	kubeClient, err := kube.NewKubeClient(logger.WithField("component", "kubeClient"))
	if err != nil {
		return errors.Wrap(err, "unable to create kube client")
	}

	analyzeService, err := kubeClient.GetServiceByLabels(AnalyzeLabelSet)
	if err != nil {
		return errors.New("failed to find analyze service")
	}

	var analyzeApiPort, analyzeServiceName string = "", ""

	for _, port := range analyzeService.Spec.Ports {
		if port.Name == "http" {
			analyzeApiPort = strconv.Itoa(int(port.Port))
			break
		}
	}

	analyzeServiceName = analyzeService.Name

	logger.Debugf("analyze service name %v, service port %v", analyzeServiceName, analyzeApiPort)

	pluginsEndpointUri := "http://" + analyzeServiceName + ":" + analyzeApiPort + "/api/v1/plugins"
	resp, err := http.Post( pluginsEndpointUri, "application/json", strings.NewReader(pluginInfo))
	if err != nil {
		return errors.Wrap(err, "failed to register plugin")
	}

	if resp.StatusCode != http.StatusCreated || resp.StatusCode != http.StatusOK {
		return errors.Errorf("failed to register plugin, status code %v", resp.StatusCode)
	}

	resp, err = http.Get(pluginsEndpointUri)
	if err != nil {
		return errors.Wrap(err, "unable to check registered plugin")
	}

	result := []*models.Plugin{}
	//TODO: generate swagger client?

	bytes, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return errors.Wrap(err, "unable to read registered plugins")
	}
	err = json.Unmarshal(bytes, result)
	if err != nil {
		return errors.Wrap(err, "unable to unmarshal registered plugins")
	}

	logger.Debugf("analyze plugins %v", string(bytes))
	for _, p := range result {
		logger.Debugf("analyze plugins %+v", p)
	}

	return nil
}

func discoverPluginApiAddress() string {
	address, exists := os.LookupEnv("PLUGIN_API_ADDRESS")
	if !exists {
		return ""
	}
	return address
}