package robot

import (
	"github.com/coreos/etcd/clientv3"
	"github.com/supergiant/robot/api"
	"github.com/supergiant/robot/pkg/logger"
)

// Config  struct represents configuration of robot service
type Config struct {
	Logging logger.Config   `mapstructure:"logging"`
	API     api.Config      `mapstructure:"api"`
	ETCD    clientv3.Config `mapstructure:"etcd"`
}

// Validate checks configuration instance for correctness
func (c *Config) Validate() error {
	return c.Logging.Validate()
}
