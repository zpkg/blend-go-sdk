package cron

import (
	"time"

	"github.com/blend/go-sdk/configutil"
)

// HistoryConfig governs job history retention in memory.
type HistoryConfig struct {
	MaxCount int           `json:"maxCount" yaml:"maxCount" env:"CRON_MAX_COUNT"`
	MaxAge   time.Duration `json:"maxAge" yaml:"maxAge" env:"CRON_MAX_AGE"`
}

// Resolve adds extra resolution steps when reading the config.
func (hc *HistoryConfig) Resolve() error {
	return configutil.AnyError(
		configutil.SetInt(&hc.MaxCount, configutil.Int(hc.MaxCount), configutil.Parse(configutil.Env("CRON_MAX_COUNT")), configutil.Int(DefaultMaxCount)),
		configutil.SetDuration(&hc.MaxAge, configutil.Duration(hc.MaxAge), configutil.Parse(configutil.Env("CRON_MAX_AGE")), configutil.Duration(DefaultMaxAge)),
	)
}

// MaxCountOrDefault returns the max count or a default.
func (hc *HistoryConfig) MaxCountOrDefault() int {
	if hc.MaxCount > 0 {
		return hc.MaxCount
	}
	return DefaultMaxCount
}

// MaxAgeOrDefault returns the max age or a default.
func (hc *HistoryConfig) MaxAgeOrDefault() time.Duration {
	if hc.MaxAge > 0 {
		return hc.MaxAge
	}
	return DefaultMaxAge
}
