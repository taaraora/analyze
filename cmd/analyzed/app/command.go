package app

import (
	"context"
	"github.com/go-openapi/loads"
	"github.com/go-openapi/strfmt"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/supergiant/analyze/pkg/analyze"
	"github.com/supergiant/analyze/pkg/api"
	"github.com/supergiant/analyze/pkg/api/handlers"
	"github.com/supergiant/analyze/pkg/api/operations"
	"github.com/supergiant/analyze/pkg/config"
	"github.com/supergiant/analyze/pkg/kube"
	"github.com/supergiant/analyze/pkg/logger"
	"github.com/supergiant/analyze/pkg/models"
	"github.com/supergiant/analyze/pkg/plugin"
	"github.com/supergiant/analyze/pkg/plugin/proto"
	"github.com/supergiant/analyze/pkg/scheduler"
	"github.com/supergiant/analyze/pkg/storage"
	"github.com/supergiant/analyze/pkg/storage/etcd"
	"os"
	"time"
)

func RunCommand(cmd *cobra.Command, _ []string) error {
	configFilePaths, err := cmd.Flags().GetStringArray("config")
	if err != nil {
		return errors.Wrap(err, "unable to get config flag value")
	}

	cfg := &analyze.Config{}

	// configFileReadError is not critical due to possibility that configuration is done by environment variables
	configFileReadError := config.ReadFromFiles(cfg, configFilePaths)

	if err = config.MergeEnv("AZ", cfg); err != nil {
		return errors.Wrap(err, "unable to merge env variables")
	}

	//TODO: try to unify APIs discovery which are hosted in k8s
	//TODO: and rewrite config population logic
	if etcdEndpoint := discoverETCDEndpoint(); etcdEndpoint != "" {
		cfg.ETCD.Endpoints = append(cfg.ETCD.Endpoints, main2.discoverETCDEndpoint())
	}

	log := logger.NewLogger(cfg.Logging).WithField("app", "analyze-core")
	mainLogger := log.WithField("component", "main")

	mainLogger.Infof("config: %+v", cfg)
	mainLogger.Infof("config file name: %s", config.UsedFileName())
	if configFileReadError != nil {
		log.Warnf("unable to read config file, %v", configFileReadError)
	}

	if err := cfg.Validate(); err != nil {
		return errors.Wrap(err, "config validation error")
	}

	kubeClient, err := kube.NewKubeClient(log.WithField("component", "kubeClient"))
	if err != nil {
		return errors.Wrap(err, "unable to create kube client")
	}

	etcdStorage, err := etcd.NewETCDStorage(cfg.ETCD, log.WithField("component", "etcdClient"))
	if err != nil {
		return errors.Wrap(err, "unable to create ETCD client")
	}

	defer etcdStorage.Close()

	scheduler := scheduler.NewScheduler(log.WithField("component", "scheduler"))
	defer scheduler.Stop()

	watchChan := etcdStorage.WatchRange(context.Background(), models.PluginPrefix)
	log.Debug("watch stated")
	pluginController := analyze.NewPluginController(
		watchChan,
		etcdStorage,
		kubeClient,
		scheduler,
		log.WithField("component", "pluginController"),
		)

	go pluginController.Loop()

	swaggerSpec, err := loads.Analyzed(api.SwaggerJSON, "2.0")
	if err != nil {
		return errors.Wrap(err, "unable to create spec analyzed document")
	}

	//TODO: add request logging middleware
	//TODO: add metrics middleware
	analyzeAPI := operations.NewAnalyzeAPI(swaggerSpec)

	analyzeAPI.GetCheckResultsHandler = handlers.NewChecksResultsHandler(
		etcdStorage,
		log.WithField("handler", "CheckResultsHandler"),
	)
	analyzeAPI.GetPluginHandler = handlers.NewPluginHandler(
		etcdStorage,
		log.WithField("handler", "PluginHandler"),
	)
	analyzeAPI.GetPluginsHandler = handlers.NewPluginsHandler(
		etcdStorage,
		log.WithField("handler", "PluginsHandler"),
	)
	analyzeAPI.RegisterPluginHandler = handlers.NewRegisterPluginHandler(
		etcdStorage,
		log.WithField("handler", "RegisterPluginHandler"),
	)
	analyzeAPI.UnregisterPluginHandler = handlers.NewUnregisterPluginHandler(
		etcdStorage,
		log.WithField("handler", "UnregisterPluginHandler"),
	)

	err = analyzeAPI.Validate()
	if err != nil {
		return errors.Wrap(err, "API configuration error")
	}

	server := api.NewServer(analyzeAPI)
	server.Port = cfg.API.ServerPort
	server.Host = cfg.API.ServerHost
	server.ConfigureAPI()

	defer server.Shutdown()

	if servingError := server.Serve(); servingError != nil {
		return errors.Wrap(servingError, "unable to serve HTTP API")
	}

	return nil
}

func discoverETCDEndpoint() string {
	etcdHost, hostExists := os.LookupEnv("ETCD_SERVICE_HOST")
	etcdPort, portExists := os.LookupEnv("ETCD_SERVICE_PORT")
	if !hostExists || !portExists {
		return ""
	}
	return etcdHost + ":" + etcdPort
}
