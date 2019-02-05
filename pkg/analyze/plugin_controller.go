package analyze

import (
	"context"
	"github.com/go-openapi/strfmt"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/pkg/errors"
	"github.com/prometheus/common/log"
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




		b, err := (&models.Plugin{
			Description: pluginInfo.Description,
			ID:          pluginInfo.Id,
			InstalledAt: strfmt.DateTime(time.Now()),
			Name:        pluginInfo.Name,
			Status:      "OK", // TODO: add status to proto, than implement plugins state which will reflect it's status
			Version:     pluginInfo.Version,
		}).MarshalBinary()
		if err != nil {
			mainLogger.Errorf("unable to load pluginClient, name: %v, error %v", pluginEntry.Name, err)
		}

		err = etcdStorage.Put(ctx, models.PluginPrefix, pluginEntry.Name, b)
		if err != nil {
			mainLogger.Errorf("unable to load pluginClient, name: %v, error %v", pluginEntry.Name, err)
		}

		check := func(pluginID string, pluginClient *plugin.Client, stor storage.Interface, logger logrus.FieldLogger) func() error {
			return func() error {
				ctx, cancel := context.WithTimeout(context.Background(), cfg.Plugin.CheckTimeout)
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
					cancel()
					return errors.Errorf("unable to marshal check result, pluginClient: %s, returned error: %s", pluginID, err)
				}

				err = stor.Put(ctx, models.CheckResultPrefix, pluginID, bytes)
				if err != nil {
					cancel()
					return errors.Errorf("unable to store check result, pluginClient: %s, returned error: %s", pluginID, err)
				}
				cancel()
				return nil
			}
		}(pluginID, pluginClient, etcdStorage, log.WithField("component", "pluginsChecks"))



		if we.Type() == storage.Added {

		}

	}
}

func (pc *PluginController) parseEvent(we storage.WatchEvent) error {
	if we.Type() == storage.Error {
		return errors.Wrap(we.Err(), "plugin watchEvent returned error")

	}

	return nil
}

func (pc *PluginController) newPluginClient(we storage.WatchEvent) error {
	pluginEntry := &models.Plugin{}
	err := pluginEntry.UnmarshalBinary(we.Payload())
	if err != nil {
		return errors.Errorf("unable to unmarshal pluginEntry entity, %s", string(we.Payload()))
	}

	// pluginEntry.ServiceLabels
	service, err := pc.kubeClient.GetService(pluginEntry.ServiceName, pluginEntry.ServiceLabels)
	if err != nil {
		return errors.Errorf("unable to find service for registered plugin, %s", string(we.Payload()))
	}

	pluginClient, err := plugin.NewClient(service.Spec.ClusterIP, cfg.Plugin.ToProtoConfig())
	if err != nil {
		return errors.Errorf("unable to instantiate pluginClient client for entity, %+v", pluginEntry)
	}

	ctx, _ := context.WithTimeout(context.Background(), cfg.Plugin.CheckTimeout)
	pluginInfo, err := pluginClient.Info(ctx, &empty.Empty{})
	if err != nil {
		return errors.Errorf("unable to load pluginInfo, name: %v, error %v", pluginEntry.Name, err)
	}
}