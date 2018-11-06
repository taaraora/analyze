package robot

import (
	"strings"

	"github.com/coreos/etcd/clientv3"
	"github.com/pkg/errors"

	"github.com/supergiant/robot/pkg/api"
	"github.com/supergiant/robot/pkg/logger"
	"github.com/supergiant/robot/pkg/plugin"
)

// Config  struct represents configuration of robot service
type Config struct {
	Logging         logger.Config   `mapstructure:"logging"`
	API             api.Config      `mapstructure:"api"`
	K8sAPIServerURI string          `mapstructure:"k8s_api_server_uri"`
	Plugin          plugin.Config   `mapstructure:"plugin"`
	ETCD            clientv3.Config `mapstructure:"etcd"`
}

// Validate checks configuration instance for correctness
func (c *Config) Validate() error {
	if err := c.Logging.Validate(); err != nil {
		return err
	}

	if err := c.Plugin.Validate(); err != nil {
		return err
	}

	if len(c.ETCD.Endpoints) == 0 {
		return errors.New("etcd endpoints where not configured")
	}

	if strings.TrimSpace(c.K8sAPIServerURI) == "" {
		return errors.New("k8s API Server uri was not configured")
	}

	return nil
}
