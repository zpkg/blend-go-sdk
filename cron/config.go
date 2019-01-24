package cron

import (
	"time"

	"github.com/blend/go-sdk/configutil"
	"github.com/blend/go-sdk/env"
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
	History HistoryConfig `json:"history" yaml:"history"`
}

// HistoryConfig governs job history retention in memory.
type HistoryConfig struct {
	MaxCount int           `json:"maxCount" yaml:"maxCount" env:"CRON_MAX_COUNT"`
	MaxAge   time.Duration `json:"maxAge" yaml:"maxAge" env:"CRON_MAX_AGE"`
}

// MaxCountOrDefault returns the max count or a default.
func (hc HistoryConfig) MaxCountOrDefault() int {
	return configutil.CoalesceInt(hc.MaxCount, DefaultMaxCount)
}

// MaxAgeOrDefault returns the max age or a default.
func (hc HistoryConfig) MaxAgeOrDefault() time.Duration {
	return configutil.CoalesceDuration(hc.MaxAge, DefaultMaxAge)
}
