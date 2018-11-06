package plugin

import (
	"context"
	"time"

	"github.com/pkg/errors"

	"github.com/supergiant/robot/pkg/plugin/proto"
)

type Config struct {
	CheckInterval time.Duration `mapstructure:"check_interval"`
	CheckTimeout  time.Duration `mapstructure:"check_timeout"`
}

type PluginsSet map[string]proto.PluginClient

// TODO: refactor and implement real pluggability
func (ps PluginsSet) Load(plugin proto.PluginClient) error {
	ctx, configureCancel := context.WithTimeout(context.Background(), time.Second*1)
	defer configureCancel()
	_, err := plugin.Configure(ctx, &proto.PluginConfig{})
	if err != nil {
		return errors.Wrap(err, "unable to configure plugin")
	}

	ctx, infoCancel := context.WithTimeout(context.Background(), time.Second*1)
	defer infoCancel()
	info, err := plugin.Info(ctx, &proto.Empty{})
	if err != nil {
		return errors.Wrap(err, "unable to get plugin info")
	}

	ps[info.Id] = plugin

	return nil
}

func (c Config) Validate() error {
	if c.CheckInterval.Nanoseconds() <= 0 {
		return errors.New("plugin check interval can't be less or equal to zero")
	}

	if c.CheckTimeout.Nanoseconds() <= 0 {
		return errors.New("plugin check timeout can't be less or equal to zero")
	}

	return nil
}
