package run

import "github.com/supergiant/robot/pkg/logger"

type Config struct {
	Logging logger.Config `mapstructure:"logging"`
}

func (c *Config) Validate() error {
	return c.Logging.Validate()
}
