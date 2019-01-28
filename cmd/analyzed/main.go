package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/supergiant/analyze/pkg/kube"

	"github.com/go-openapi/loads"
	"github.com/go-openapi/strfmt"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/supergiant/analyze"
	"github.com/supergiant/analyze/pkg/api"
	"github.com/supergiant/analyze/pkg/api/handlers"
	"github.com/supergiant/analyze/pkg/api/operations"
	"github.com/supergiant/analyze/pkg/config"
	"github.com/supergiant/analyze/pkg/logger"
	"github.com/supergiant/analyze/pkg/models"
	"github.com/supergiant/analyze/pkg/plugin"
	"github.com/supergiant/analyze/pkg/plugin/proto"
	"github.com/supergiant/analyze/pkg/scheduler"
	"github.com/supergiant/analyze/pkg/storage"
	"github.com/supergiant/analyze/pkg/storage/etcd"
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

	cfg := &analyze.Config{}

	// configFileReadError is not critical due to possibility that configuration is done by environment variables
	configFileReadError := config.ReadFromFiles(cfg, configFilePaths)

	if err = config.MergeEnv("AZ", cfg); err != nil {
		return errors.Wrap(err, "unable to merge env variables")
	}

	//TODO: try to unify APIs discovery which are hosted in k8s
	//TODO: and rewrite config population logic
	if etcdEndpoint := discoverETCDEndpoint(); etcdEndpoint != "" {
		cfg.ETCD.Endpoints = append(cfg.ETCD.Endpoints, discoverETCDEndpoint())
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

	etcdStorage, err := etcd.NewETCDStorage(cfg.ETCD)
	if err != nil {
		return errors.Wrap(err, "unable to create ETCD client")
	}

	defer etcdStorage.Close()
	//load registered plugins
	rawPluginsEntries, err := etcdStorage.GetAll(context.Background(), models.PluginPrefix)
	if err != nil {
		return errors.Wrap(err, "unable to read plugins registry from ETCD")
	}

	plugins := make(plugin.PluginsSet)

	for _, rawPluginEntry := range rawPluginsEntries {
		pluginEntry := &models.Plugin{}
		err := pluginEntry.UnmarshalBinary(rawPluginEntry)
		if err != nil {
			log.Warnf("unable to unmarshal pluginEntry entity, %s", string(rawPluginEntry))
			continue
		}

		// pluginEntry.ServiceLabels
		service, err := kubeClient.GetService(pluginEntry.ServiceName, pluginEntry.ServiceLabels)
		if err != nil {
			log.Warnf("unable to find service for registered plugin, %s", string(rawPluginEntry))
			continue
		}

		pluginClient, err := plugin.NewClient(service.Spec.ClusterIP, cfg.Plugin.ToProtoConfig())
		if err != nil {
			log.Warnf("unable to instantiate pluginClient client for entity, %+v", pluginEntry)
			continue
		}
		plugins[pluginEntry.ID] = pluginClient
	}

	for pluginName, pluginClient := range plugins {
		ctx, _ := context.WithTimeout(context.Background(), cfg.Plugin.CheckTimeout)
		pluginInfo, err := pluginClient.Info(ctx, &empty.Empty{})
		if err != nil {
			mainLogger.Errorf("unable to load pluginClient, name: %v, error %v", pluginName, err)
		}

		b, err := (&models.Plugin{
			Description: pluginInfo.Description,
			ID:          pluginInfo.Id,
			InstalledAt: strfmt.DateTime(time.Now()),
			Name:        pluginInfo.Name,
			Status:      "OK", // TODO: add status to proto, than implement plugins state which will reflect it's status
			Version:     pluginInfo.Version,
		}).MarshalBinary()
		if err != nil {
			mainLogger.Errorf("unable to load pluginClient, name: %v, error %v", pluginName, err)
		}

		err = etcdStorage.Put(ctx, models.PluginPrefix, pluginName, b)
		if err != nil {
			mainLogger.Errorf("unable to load pluginClient, name: %v, error %v", pluginName, err)
		}
	}

	check := func(p plugin.PluginsSet, stor storage.Interface, logger logrus.FieldLogger) func() {
		return func() {
			for pluginID, pluginClient := range p {
				ctx, cancel := context.WithTimeout(context.Background(), cfg.Plugin.CheckTimeout)
				checkResponse, err := pluginClient.Check(ctx, &proto.CheckRequest{})
				if err != nil {
					logger.Errorf("unable to execute check for pluginClient: %s, error: %v", pluginID, err)
					cancel()
					continue
				}
				if checkResponse.Error != "" {
					logger.Errorf("pluginClient: %s, returned error: %s", pluginID, checkResponse.Error)
					cancel()
					continue
				}

				if checkResponse.Result == nil {
					logger.Errorf("pluginClient: %s, returned nil Result", pluginID)
					cancel()
					continue
				}

				r := checkResponse.Result

				var actions = []*models.PluginAction{}
				for _, action := range r.Actions {
					actions = append(actions, &models.PluginAction{
						Description: action.Description,
						Name:        action.Name,
						ID:          action.ActionId,
					})
				}
				var currentTime = time.Now()
				checkResult := models.CheckResult{
					CheckStatus:     r.GetStatus().String(),
					CompletedAt:     strfmt.DateTime(currentTime),
					Description:     string(r.GetDescription().Value),
					ExecutionStatus: r.GetExecutionStatus(),
					ID:              r.GetName(),
					Name:            r.GetName(),
					PossibleActions: actions,
				}

				bytes, err := checkResult.MarshalBinary()
				if err != nil {
					logger.Errorf("unable to marshal check result, pluginClient: %s, returned error: %s", pluginID, err)
					cancel()
					continue
				}

				err = stor.Put(ctx, models.CheckResultPrefix, pluginID, bytes)
				if err != nil {
					logger.Errorf("unable to store check result, pluginClient: %s, returned error: %s", pluginID, err)
					cancel()
				}

				cancel()
			}
		}
	}(plugins, etcdStorage, log.WithField("component", "pluginsChecks"))

	scheduler := scheduler.NewScheduler(cfg.Plugin.CheckInterval, check)
	defer scheduler.Stop()

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

func logEnvs(logger logrus.FieldLogger) {
	for _, pair := range os.Environ() {
		logger.Warnf("%s", pair)
	}
}

func discoverETCDEndpoint() string {
	etcdHost, hostExists := os.LookupEnv("ETCD_SERVICE_HOST")
	etcdPort, portExists := os.LookupEnv("ETCD_SERVICE_PORT")
	if !hostExists || !portExists {
		return ""
	}
	return etcdHost + ":" + etcdPort
}
