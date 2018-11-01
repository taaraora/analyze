package main

import (
	"github.com/sirupsen/logrus"
	"log"
	"os"

	"github.com/supergiant/robot"
	"github.com/supergiant/robot/pkg/api"
	"github.com/supergiant/robot/pkg/api/handlers"
	"github.com/supergiant/robot/pkg/api/operations"
	"github.com/supergiant/robot/pkg/config"
	"github.com/supergiant/robot/pkg/logger"
	"github.com/supergiant/robot/pkg/storage/etcd"

	"github.com/go-openapi/loads"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func main() {

	command := &cobra.Command{
		Use:          "analyzed",
		Short:        "analyze service checks K8s cluster by means of installed plugins and gives recommendations",
		RunE:         runCommand,
		SilenceUsage: true,
	}

	command.PersistentFlags().StringArrayP(
		"config",
		"c",
		[]string{"./config.yaml", "/etc/analyzed/config.yaml", "$HOME/.analyzed/config.yaml"},
		"config file path")

	if err := command.Execute(); err != nil {
		log.Fatalf("\n%v\n", err)
	}
}

func runCommand(cmd *cobra.Command, _ []string) error {
	configFilePaths, err := cmd.Flags().GetStringArray("config")
	if err != nil {
		return errors.Wrap(err, "unable to get config flag value")
	}

	cfg := &robot.Config{}

	// configFileReadError is not critical due to possibility that configuration is done by environment variables
	configFileReadError := config.ReadFromFiles(cfg, configFilePaths)

	if err = config.MergeEnv("AZ", cfg); err != nil {
		return errors.Wrap(err, "unable to merge env variables")
	}

	log := logger.NewLogger(cfg.Logging).WithField("app", "robot")
	mainLogger := log.WithField("component", "main")

	mainLogger.Infof("config: %+v", cfg)
	mainLogger.Infof("config file name: %s", config.UsedFileName())
	if configFileReadError != nil {
		mainLogger.Warnf("unable to read config file, %v", configFileReadError)
	}

	if err := cfg.Validate(); err != nil {
		return errors.Wrap(err, "config validation error")
	}

	storage, err := etcd.NewETCDStorage(cfg.ETCD)
	if err != nil {
		return errors.Wrap(err, "unable to create ETCD client")
	}

	defer storage.Close()

	swaggerSpec, err := loads.Analyzed(api.SwaggerJSON, "2.0")
	if err != nil {
		return errors.Wrap(err, "unable to create spec analyzed document")
	}

	analyzeAPI := operations.NewAnalyzeAPI(swaggerSpec)
	analyzeAPI.GetRecommendationPluginsHandler = handlers.NewRecommendationPluginsHandler(storage)
	analyzeAPI.GetCheckResultsHandler = handlers.NewCheckResultsHandler(storage)

	server := api.NewServer(analyzeAPI)
	server.Port = cfg.API.ServerPort
	server.Host = cfg.API.ServerHost
	server.ConfigureAPI()

	defer server.Shutdown()

	loadEnvs(mainLogger)

	if servingError := server.Serve(); servingError != nil {
		return errors.Wrap(servingError, "unable to serve HTTP API")
	}

	return nil
}

func loadEnvs(logger logrus.FieldLogger)  {
	kubeHost, _ := os.LookupEnv("KUBERNETES_SERVICE_HOST")
	logger.Warnf("%s", kubeHost)
	kubePort, _ := os.LookupEnv("KUBERNETES_SERVICE_PORT")
	logger.Warnf("%s", kubePort)
}