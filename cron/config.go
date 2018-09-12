package cron

import (
	"time"

	"github.com/blend/go-sdk/env"
	"github.com/blend/go-sdk/util"
)

// NewConfigFromEnv creates a new config from the environment.
func NewConfigFromEnv() (*Config, error) {
	var cfg Config
	if err := env.Env().ReadInto(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

// MustNewConfigFromEnv returns a new config set from environment variables,
// it will panic if there is an error.
func MustNewConfigFromEnv() *Config {
	cfg, err := NewConfigFromEnv()
	if err != nil {
		panic(err)
	}
	return cfg
}

// Config is the config object.
type Config struct {
	HeartbeatInterval time.Duration `json:"heartbeatInterval" yaml:"heartbeatInterval" env:"CRON_HEARTBEAT_INTERVAL"`
}

// GetHeartbeatInterval gets a property or a default.
func (c Config) GetHeartbeatInterval(inherited ...time.Duration) time.Duration {
	return util.Coalesce.Duration(c.HeartbeatInterval, DefaultHeartbeatInterval, inherited...)
}
