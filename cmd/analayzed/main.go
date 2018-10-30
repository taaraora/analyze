package main

import (
	"github.com/supergiant/robot/pkg/api/handlers"
	"log"

	"github.com/supergiant/robot/pkg/api"
	"github.com/supergiant/robot/pkg/api/operations"

	"github.com/go-openapi/loads"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/supergiant/robot"
	"github.com/supergiant/robot/pkg/config"
	"github.com/supergiant/robot/pkg/logger"
	"github.com/supergiant/robot/pkg/storage/etcd"
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
		[]string{"./analyzed.yaml", "/etc/analyzed/analyzed.yaml", "$HOME/analyzed.yaml"},
		"config file path")

	if err := command.Execute(); err != nil {
		log.Fatalf("\n%v\n", err)
	}
}

func runCommand(cmd *cobra.Command, _ []string) error {
	configFilePaths, flagReadError := cmd.Flags().GetStringArray("config")
	if flagReadError != nil {
		return errors.Wrap(flagReadError, "unable to get config flag value")
	}

	cfg := &robot.Config{}

	// configFileReadError is not critical due to possibility that configuration is done by environment variables
	configFileReadError := config.ReadFromFiles(cfg, configFilePaths)

	if envVariablesReadError := config.MergeEnv("RK", cfg); envVariablesReadError != nil {
		return errors.Wrap(envVariablesReadError, "unable to merge env variables")
	}

	log := logger.NewLogger(cfg.Logging).WithField("app", "robot")
	mainLogger := log.WithField("component", "main")

	mainLogger.Infof("config: %+v", cfg)
	mainLogger.Infof("config file name: %s", config.UsedFileName())
	if configFileReadError != nil {
		mainLogger.Warnf("unable to read config file, %v", configFileReadError)
	}

	if cfgValidationError := cfg.Validate(); cfgValidationError != nil {
		return errors.Wrap(cfgValidationError, "config validation error")
	}

	swaggerSpec, specDocumentCreationError := loads.Analyzed(api.SwaggerJSON, "2.0")
	if specDocumentCreationError != nil {
		return errors.Wrap(specDocumentCreationError, "unable to create spec analyzed document")
	}

	analyzeAPI := operations.NewAnalyzeAPI(swaggerSpec)
	server := api.NewServer(analyzeAPI)
	defer server.Shutdown()
	server.Port = cfg.API.ServerPort
	server.Host = cfg.API.ServerHost

	storage, err := etcd.NewETCDStorage(cfg.ETCD)
	if err != nil {
		errors.Wrap(specDocumentCreationError, "unable to create ETCD client")
	}
	defer storage.Close()

	analyzeAPI.GetRecommendationPluginsHandler = handlers.NewRecommendationPluginsHandler(storage)
	analyzeAPI.GetCheckResultsHandler = handlers.NewCheckResultsHandler(storage)
	server.ConfigureAPI()

	if err := server.Serve(); err != nil {
		mainLogger.Fatal(err)
	}

	return nil
}
