package robot

import "github.com/supergiant/robot/pkg/logger"

// Config  struct represents configuration of robot service
type Config struct {
	Logging logger.Config `mapstructure:"logging"`
}

// Validate checks configuration instance for correctness
func (c *Config) Validate() error {
	return c.Logging.Validate()
}
