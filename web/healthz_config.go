package web

import (
	"time"

	"github.com/blend/go-sdk/util"
)

// HealthzConfig is the healthz config.
type HealthzConfig struct {
	BindAddr      string        `json:"bindAddr" yaml:"bindAddr" env:"HEALTHZ_BIND_ADDR"`
	GracePeriod   time.Duration `json:"gracePeriod" yaml:"gracePeriod"`
	RecoverPanics *bool         `json:"recoverPanics" yaml:"recoverPanics"`
}

// GetBindAddr gets the bind address.
func (hzc HealthzConfig) GetBindAddr(defaults ...string) string {
	return util.Coalesce.String(hzc.BindAddr, DefaultHealthzBindAddr, defaults...)
}

// GetGracePeriod gets a grace period or a default.
func (hzc HealthzConfig) GetGracePeriod(defaults ...time.Duration) time.Duration {
	return util.Coalesce.Duration(hzc.GracePeriod, DefaultShutdownGracePeriod, defaults...)
}

// GetRecoverPanics gets recover panics or a default.
func (hzc HealthzConfig) GetRecoverPanics(defaults ...bool) bool {
	return util.Coalesce.Bool(hzc.RecoverPanics, DefaultRecoverPanics, defaults...)
}
