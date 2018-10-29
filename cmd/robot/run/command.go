package run

import (
	"github.com/go-openapi/loads"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/supergiant/robot"
	"github.com/supergiant/robot/pkg/config"
	"github.com/supergiant/robot/pkg/logger"
	"github.com/supergiant/robot/swagger/gen/restapi"
	"github.com/supergiant/robot/swagger/gen/restapi/operations"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "run",
		Short: "runs robot as a service hosted in k8s cluster",
		RunE:  run,
	}

	return cmd
}

func run(cmd *cobra.Command, _ []string) error {
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

	swaggerSpec, specDocumentCreationError := loads.Analyzed(restapi.SwaggerJSON, "2.0")
	if specDocumentCreationError != nil {
		return errors.Wrap(specDocumentCreationError, "unable to create spec analyzed document")
	}

	api := operations.NewRobotAPI(swaggerSpec)
	server := restapi.NewServer(api)
	defer server.Shutdown()
	server.Port = 9091
	server.ConfigureAPI()

	if err := server.Serve(); err != nil {
		mainLogger.Fatal(err)
	}

	return nil
}
