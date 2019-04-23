package plugin

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/supergiant/analyze/pkg/storage/etcd"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/supergiant/analyze/pkg/kube"
	"github.com/supergiant/analyze/pkg/models"
	"github.com/supergiant/analyze/pkg/plugin/proto"
	"github.com/supergiant/analyze/pkg/proxy"
	"github.com/supergiant/analyze/pkg/scheduler"
	"github.com/supergiant/analyze/pkg/storage"
)

type Controller struct {
	storeConfig        *etcd.Config
	checkJob           func(string, *Client, storage.Interface) func() error
	pluginChangeEvents <-chan storage.WatchEvent
	stor               storage.Interface
	kubeClient         kube.Interface
	logger             logrus.FieldLogger
	scheduler          scheduler.Interface
	proxySet           *proxy.Set
	pluginClients      map[string]*Client
	stop               chan struct{}
}

func NewPluginController(
	cfg *etcd.Config,
	checkJob func(string, *Client, storage.Interface) func() error,
	events <-chan storage.WatchEvent,
	stor storage.Interface,
	kubeClient kube.Interface,
	scheduler scheduler.Interface,
	set *proxy.Set,
	logger logrus.FieldLogger,
) *Controller {
	pc := &Controller{
		storeConfig:        cfg,
		checkJob:           checkJob,
		pluginChangeEvents: events,
		stor:               stor,
		kubeClient:         kubeClient,
		scheduler:          scheduler,
		proxySet:           set,
		pluginClients:      make(map[string]*Client),
		logger:             logger,
		stop:               make(chan struct{}),
	}

	go pc.handlePluginChange()
	return pc
}

func (pc *Controller) handlePluginChange() {
	for {
		select {
		case <-pc.stop:
			return
		case we := <-pc.pluginChangeEvents:
			if err := pc.handleEvent(we); err != nil {
				pc.logger.Errorf("unable to handle watch event, err: %+v", err)
			}
		}
	}
}

// TODO: maybe split config and plugin updates?
func (pc *Controller) handleEvent(we storage.WatchEvent) error {
	if we == nil {
		return errors.New("plugin watchEvent is nil")
	}
	if we.Type() == storage.Error {
		return errors.Wrap(we.Err(), "plugin watchEvent returned error")
	}
	var err error

	pluginEntry := &models.Plugin{}
	err = pluginEntry.UnmarshalBinary(we.Payload()) // TODO: handle DELETE logic, it requires storage interface update
	if err != nil {
		pc.logger.Errorf("we payload '%s'", string(we.Payload()))
		return errors.Errorf("unable to unmarshal pluginEntry entity err: %v", err)
	}

	if we.Type() == storage.Deleted {
		return pc.unregisterPlugin(pluginEntry)
	}

	if we.Type() == storage.Added {
		return pc.registerPlugin(pluginEntry)
	}

	if we.Type() == storage.Modified {
		err = pc.unregisterPlugin(pluginEntry)
		if err != nil {
			pc.logger.Errorf("unable to unregister modified plugin, error: %v", err)
		}
		return pc.registerPlugin(pluginEntry)
	}

	return nil
}

// newPluginClient returns plugin client instance which is configured and ready for work
func (pc *Controller) newPluginClient(pluginEntry *models.Plugin) (*Client, error) {

	// pluginEntry.ServiceLabels
	ss := strings.Split(pluginEntry.ServiceEndpoint, ":")
	if len(ss) < 2 {
		return nil, errors.Errorf("unable get service name, %+v", pluginEntry.ServiceEndpoint)
	}
	// TODO: empty namespace means that serviceaccount of analyze pod
	//  need to have access to all namespaces to list services?
	service, err := pc.kubeClient.GetService(ss[0], "", pluginEntry.ServiceLabels)
	if err != nil {
		return nil, errors.Errorf("unable to find service for registered plugin, %+v", pluginEntry)
	}

	var servicePort string
	for _, port := range service.Spec.Ports {
		if port.Name == "grpc" {
			servicePort = strconv.Itoa(int(port.Port))
		}
	}

	if servicePort == "" {
		return nil, errors.Errorf(
			"unable to find service  port for plugin, %+v, list of ports: %+v",
			pluginEntry,
			service.Spec.Ports,
		)
	}

	pluginClient, err := NewClient(service.Spec.ClusterIP + ":" + servicePort)
	if err != nil {
		return nil, errors.Errorf("unable to instantiate pluginClient client for entity, %+v", pluginEntry)
	}

	return pluginClient, nil
}

func (pc *Controller) unregisterPlugin(pluginEntry *models.Plugin) error {
	pluginClient, exists := pc.pluginClients[pluginEntry.ID]
	if !exists {
		return errors.Errorf("unable to find pluginClient, name: %v", pluginEntry.ID)
	}

	delete(pc.pluginClients, pluginEntry.ID)

	err := pc.scheduler.RemoveJob(pluginEntry.ID)
	if err != nil {
		return errors.Errorf("unable to find plugin job in scheduler, id: %v ", pluginEntry.ID)
	}

	err = pluginClient.Close()
	if err != nil {
		return errors.Errorf("unable to close pluginClient, id: %v, err: %+v", pluginEntry.ID, err)
	}

	pc.proxySet.RemoveProxy(pluginEntry.ID)

	return nil
}

func (pc *Controller) registerPlugin(pluginEntry *models.Plugin) error {
	var pluginConfig *proto.PluginConfig
	pluginClient, err := pc.newPluginClient(pluginEntry)
	if err != nil {
		return errors.Wrap(err, "can't create plugin client for watchEvent")
	}

	// TODO: make it configurable
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*120)
	defer cancel()
	pluginInfo, err := pluginClient.Info(ctx, &empty.Empty{})
	if err != nil {
		return errors.Errorf(
			"unable to load pluginInfo, name: %v, target: %s, error %v",
			pluginEntry.ID,
			pluginClient.GetConnTarget(),
			err,
		)
	}
	pc.logger.Infof("plugin info loaded %+v", *pluginInfo)

	rawConfig, err := pc.stor.Get(ctx, storage.PluginConfigPrefix, pluginInfo.Id)
	pc.logger.Infof("plugin raw config %s, error: %+v", rawConfig, err)
	if err != nil && err != storage.ErrNotFound {
		return errors.Errorf("unable to load plugin config, id: %v, error %v", pluginEntry.ID, err)
	}

	if err == storage.ErrNotFound && pluginInfo.DefaultConfig == nil {
		return errors.Errorf("unable to register plugin, without default config, id: %v", pluginEntry.ID)
	}

	// when plugin is just installed we need to take default config from plugin and save it
	if err == storage.ErrNotFound {
		pluginConfig = pluginInfo.DefaultConfig
		pluginConfig.EtcdEndpoints = pc.storeConfig.Endpoints
		pc.logger.Infof(
			"plugin %s config not found, etcd endpoints: %+v, default config %+v",
			pluginInfo.Id,
			pluginConfig.EtcdEndpoints,
			pluginInfo.DefaultConfig,
		)

		var pluginSpecificConfig = make(map[string]interface{})
		if pluginInfo.DefaultConfig.PluginSpecificConfig != nil {
			marshalErr := json.Unmarshal(pluginInfo.DefaultConfig.PluginSpecificConfig, &pluginSpecificConfig)
			if marshalErr != nil {
				return errors.Errorf("unable to marshal plugin default config, id: %v, error %v", pluginEntry.ID, marshalErr)
			}
			pc.logger.Infof(
				"plugin %s specific config unmarshalled %+v",
				pluginInfo.Id,
				pluginSpecificConfig,
			)
		}

		pluginConfigEntry := &models.PluginConfig{
			ExecutionInterval:    pluginInfo.DefaultConfig.ExecutionInterval.Seconds,
			EtcdEndpoints:        pc.storeConfig.Endpoints,
			PluginSpecificConfig: &pluginSpecificConfig,
		}

		b, marshalErr := pluginConfigEntry.MarshalBinary()
		if marshalErr != nil {
			return errors.Errorf("unable to marshal plugin default config, id: %v, error %v", pluginEntry.ID, marshalErr)
		}
		putError := pc.stor.Put(ctx, storage.PluginConfigPrefix, pluginInfo.Id, msg(b))
		if putError != nil {
			return errors.Errorf("unable to store plugin default config, id: %v, error %v", pluginEntry.ID, putError)
		}
	}

	// true means plugin config has been read from storage and we need to apply it (case when analyze core was rebooted)
	if pluginConfig == nil {
		pluginConfigEntry := &models.PluginConfig{}
		err = pluginConfigEntry.UnmarshalBinary(rawConfig.Payload())
		if err != nil {
			return errors.Errorf("unable to unmarshal plugin config, id: %v, error %v", pluginEntry.ID, err)
		}

		bytes, err := json.Marshal(pluginConfigEntry.PluginSpecificConfig)
		if err != nil {
			return errors.Errorf("unable to marshal plugin config, id: %v, error %v", pluginEntry.ID, err)
		}

		// TODO: populate other properties
		pluginConfig = &proto.PluginConfig{
			ExecutionInterval:    ptypes.DurationProto(time.Second * time.Duration(pluginConfigEntry.ExecutionInterval)),
			EtcdEndpoints:        pc.storeConfig.Endpoints,
			PluginSpecificConfig: bytes,
		}
	}

	pc.pluginClients[pluginInfo.Id] = pluginClient

	_, err = pluginClient.Configure(ctx, pluginConfig)
	if err != nil {
		return errors.Wrap(err, "unable to configure plugin")
	}

	interval, err := ptypes.Duration(pluginConfig.ExecutionInterval)
	if err != nil {
		return errors.Wrap(err, "unable to parse execution interval for plugin")
	}

	err = pc.scheduler.ScheduleJob(pluginInfo.Id, interval, pc.checkJob(pluginInfo.Id, pluginClient, pc.stor))
	if err != nil {
		return errors.Wrap(err, "unable to schedule job for plugin")
	}

	err = pc.proxySet.SetProxy(pluginEntry)
	if err != nil {
		return errors.Wrap(err, "unable to register proxy for plugin")
	}

	return nil
}

func (pc *Controller) Stop() {
	pc.stop <- struct{}{}
}

type msg []byte

func (d msg) Payload() []byte {
	return d
}
