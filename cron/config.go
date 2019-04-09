package cron

import (
	"time"

	"github.com/blend/go-sdk/configutil"
)

// Config governs job history retention in memory.
type Config struct {
	HistoryMaxCount int           `json:"historyMaxCount" yaml:"historyMaxCount" env:"CRON_HISTORY_MAX_COUNT"`
	HistoryMaxAge   time.Duration `json:"historyMaxAge" yaml:"historyMaxAge" env:"CRON_HISTORY_MAX_AGE"`
}

// Resolve adds extra resolution steps when reading the config.
func (hc Config) Resolve() error {
	return configutil.AnyError(
		configutil.SetInt(&hc.HistoryMaxCount, configutil.Int(hc.HistoryMaxCount), configutil.Parse(configutil.Env("CRON_HISTORY_MAX_COUNT")), configutil.Int(DefaultHistoryMaxCount)),
		configutil.SetDuration(&hc.HistoryMaxAge, configutil.Duration(hc.HistoryMaxAge), configutil.Parse(configutil.Env("CRON_HISTORY_MAX_AGE")), configutil.Duration(DefaultHistoryMaxAge)),
	)
}

// HistoryMaxCountOrDefault returns the max count or a default.
func (hc Config) HistoryMaxCountOrDefault() int {
	if hc.HistoryMaxCount > 0 {
		return hc.HistoryMaxCount
	}
	return DefaultHistoryMaxCount
}

// HistoryMaxAgeOrDefault returns the max age or a default.
func (hc Config) HistoryMaxAgeOrDefault() time.Duration {
	if hc.HistoryMaxAge > 0 {
		return hc.HistoryMaxAge
	}
	return DefaultHistoryMaxAge
}
