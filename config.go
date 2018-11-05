package robot

import (
	"strings"
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/pkg/errors"

	"github.com/supergiant/robot/pkg/api"
	"github.com/supergiant/robot/pkg/logger"
)

// Config  struct represents configuration of robot service
type Config struct {
	Logging         logger.Config   `mapstructure:"logging"`
	API             api.Config      `mapstructure:"api"`
	K8sAPIServerURI string          `mapstructure:"k8s_api_server_uri"`
	ETCD            clientv3.Config `mapstructure:"etcd"`
	CheckInterval   time.Duration   `mapstructure:"check_interval"`
}

// Validate checks configuration instance for correctness
func (c *Config) Validate() error {
	if err := c.Logging.Validate(); err != nil {
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
