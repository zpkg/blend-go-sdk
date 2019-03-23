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

// MaxCountOrDefault returns the max count or a default.
func (hc HistoryConfig) MaxCountOrDefault() int {
	return configutil.CoalesceInt(hc.MaxCount, DefaultMaxCount)
}

// MaxAgeOrDefault returns the max age or a default.
func (hc HistoryConfig) MaxAgeOrDefault() time.Duration {
	return configutil.CoalesceDuration(hc.MaxAge, DefaultMaxAge)
}
