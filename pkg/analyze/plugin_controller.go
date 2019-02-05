package analyze

import (
	"context"
	"github.com/go-openapi/strfmt"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/supergiant/analyze/pkg/kube"
	"github.com/supergiant/analyze/pkg/models"
	"github.com/supergiant/analyze/pkg/plugin"
	"github.com/supergiant/analyze/pkg/plugin/proto"
	"github.com/supergiant/analyze/pkg/storage"
	"time"
)

type PluginController struct {
	events <-chan storage.WatchEvent
	stor storage.Interface
	kubeClient kube.Interface
	logger logrus.FieldLogger
	pluginClients map[string]*plugin.Client
}

func NewPluginController(
	events <-chan storage.WatchEvent,
	stor storage.Interface,
	kubeClient kube.Interface,
	logger logrus.FieldLogger,
	) *PluginController {
	return &PluginController{
		events: events,
		stor:stor,
		kubeClient:kubeClient,
		logger:logger,
	}
}

func (pc *PluginController) Loop(){
	for we := range pc.events {

		if err := pc.parseEvent(we); err != nil {
			pc.logger.Errorf("unable to parse watch event")
		}

	}
}

// TODO: maybe split config and plugin updates?
func (pc *PluginController) parseEvent(we storage.WatchEvent) error {
	if we.Type() == storage.Error {
		return errors.Wrap(we.Err(), "plugin watchEvent returned error")
	}
	var err error
	var pluginClient *plugin.Client
	var pluginConfig *proto.PluginConfig

	pluginEntry := &models.Plugin{}
	err = pluginEntry.UnmarshalBinary(we.Payload())
	if err != nil {
		return errors.Errorf("unable to unmarshal pluginEntry entity, %s", string(we.Payload()))
	}

	if we.Type() == storage.Deleted {
		var exists bool
		pluginClient, exists = pc.pluginClients[pluginEntry.ID]
		if !exists {
			return errors.Errorf("unable to find pluginClient, name: %v", pluginEntry.Name)
		}
		err = pluginClient.Close()
		if err != nil {
			return errors.Errorf("unable to find pluginClient, name: %v", pluginEntry.Name)
		}
		delete(pc.pluginClients, pluginEntry.ID)
		return nil
	}

	if we.Type() == storage.Added {
 		pluginClient, err = pc.newPluginClient(pluginEntry)
 		if err != nil {
			return errors.Wrap(we.Err(), "can't create plugin client for watchEvent")
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
		pluginInfo, err := pluginClient.Info(ctx, &empty.Empty{})
		defer cancel()
		if err != nil {
			return errors.Errorf("unable to load pluginInfo, name: %v, error %v", pluginEntry.Name, err)
		}

		pluginConfig = pluginInfo.DefaultConfig
		pc.pluginClients[pluginInfo.Id] = pluginClient
	}

	if we.Type() == storage.Modified {
		var exists bool
		pluginClient, exists = pc.pluginClients[pluginEntry.ID]
		if !exists {
			return errors.Errorf("unable to find pluginClient, name: %v", pluginEntry.Name)
		}
	}


	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()
	_, err = pluginClient.Configure(ctx, pluginConfig)
	if err != nil {
		return errors.Wrap(err, "unable to configure plugin")
	}

	return nil
}

// newPluginClient returns plugin client instance which is configured and ready for work
func (pc *PluginController) newPluginClient(pluginEntry *models.Plugin) (*plugin.Client, error) {

	// pluginEntry.ServiceLabels
	service, err := pc.kubeClient.GetService(pluginEntry.ServiceName, pluginEntry.ServiceLabels)
	if err != nil {
		return nil, errors.Errorf("unable to find service for registered plugin, %+v", pluginEntry)
	}

	pluginClient, err := plugin.NewClient(service.Spec.ClusterIP)
	if err != nil {
		return nil, errors.Errorf("unable to instantiate pluginClient client for entity, %+v", pluginEntry)
	}

	return pluginClient, nil
}

func (pc *PluginController) check(pluginID string, pluginClient *plugin.Client) func()error {
	return func() error {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		checkResponse, err := pluginClient.Check(ctx, &proto.CheckRequest{})
		if err != nil {
			cancel()
			return errors.Errorf("unable to execute check for pluginClient: %s, error: %v", pluginID, err)
		}
		if checkResponse.Error != "" {
			cancel()
			return errors.Errorf("pluginClient: %s, returned error: %s", pluginID, checkResponse.Error)
		}

		if checkResponse.Result == nil {
			cancel()
			return errors.Errorf("pluginClient: %s, returned nil Result", pluginID)
		}

		r := checkResponse.Result

		var currentTime = time.Now()
		checkResult := models.CheckResult{
			CheckStatus:     r.GetStatus().String(),
			CompletedAt:     strfmt.DateTime(currentTime),
			Description:     string(r.GetDescription().Value),
			ExecutionStatus: r.GetExecutionStatus(),
			ID:              r.GetName(),
			Name:            r.GetName(),
		}

		bytes, err := checkResult.MarshalBinary()
		if err != nil {
			cancel()
			return errors.Errorf("unable to marshal check result, pluginClient: %s, returned error: %s", pluginID, err)
		}

		err = pc.stor.Put(ctx, models.CheckResultPrefix, pluginID, bytes)
		if err != nil {
			cancel()
			return errors.Errorf("unable to store check result, pluginClient: %s, returned error: %s", pluginID, err)
		}
		cancel()
		return nil
	}
}