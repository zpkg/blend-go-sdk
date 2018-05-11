package secrets

import (
	"time"

	"github.com/blend/go-sdk/env"
	"github.com/blend/go-sdk/util"
)

const (
	// DefaultAddr is the default addr.
	DefaultAddr = "http://127.0.0.1:8200"

	// DefaultTimeout is the default timeout.
	DefaultTimeout = time.Second
)

// NewConfigFromEnv returns a config populated by the env.
func NewConfigFromEnv() *Config {
	var cfg Config
	if err := env.Env().ReadInto(&cfg); err != nil {
		panic(err)
	}
	return &cfg
}

// Config is the secrets config object.
type Config struct {
	Addr    string        `json:"addr" yaml:"addr" env:"VAULT_ADDR"`
	Token   string        `json:"token" yaml:"token" env:"VAULT_TOKEN"`
	Timeout time.Duration `json:"timeout" yaml:"timeout"`
}

// GetAddr returns the client addr.
func (c Config) GetAddr(inherited ...string) string {
	return util.Coalesce.String(c.Addr, DefaultAddr, inherited...)
}

// GetToken returns the client token.
func (c Config) GetToken(inherited ...string) string {
	return util.Coalesce.String(c.Token, "", inherited...)
}

// GetTimeout returns the client timeout.
func (c Config) GetTimeout(inherited ...time.Duration) time.Duration {
	return util.Coalesce.Duration(c.Timeout, DefaultTimeout, inherited...)
}
