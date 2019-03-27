package secrets

import (
	"net/url"
	"time"

	"github.com/blend/go-sdk/configutil"
	"github.com/blend/go-sdk/env"
)

// EnvVars
const (
	EnvVarVaultAddr  = "VAULT_ADDR"
	EnvVarVaultToken = "VAULT_TOKEN"
)

// MustNewConfigFromEnv returns a config set from the env, and panics on error.
func MustNewConfigFromEnv() (cfg *Config) {
	var err error
	if cfg, err = NewConfigFromEnv(); err != nil {
		panic(err)
	}
	return
}

// NewConfigFromEnv returns a config populated by the env.
func NewConfigFromEnv() (*Config, error) {
	var cfg Config
	if err := env.Env().ReadInto(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

// Config is the secrets config object.
type Config struct {
	// Addr is the remote address of the secret store.
	Addr string `json:"addr" yaml:"addr" env:"VAULT_ADDR"`
	// Token is the authentication token used to talk to the secret store.
	Token string `json:"token" yaml:"token" env:"VAULT_TOKEN"`
	// Mount is the default mount path, it prefixes any keys.
	Mount string `json:"mount" yaml:"mount"`
	// Timeout is the dial timeout for requests to the secrets store.
	Timeout time.Duration `json:"timeout" yaml:"timeout"`
	// RootCAs is a list of certificate authority paths.
	RootCAs []string `json:"rootCAs" yaml:"rootCAs" env:"VAULT_CACERT,csv"`
	// ServicePath is the path that service secrets live under
	ServicePath string `json:"servicePath" yaml:"servicePath" env:"SECRETS_SERVICE_PATH"`
}

// IsZero returns if the config is set or not.
func (c Config) IsZero() bool {
	return len(c.Token) == 0
}

// Resolve reads the environment into the config on configutil.Read(...)
func (c *Config) Resolve() error {
	return env.Env().ReadInto(c)
}

// AddrOrDefault returns the client addr.
func (c Config) AddrOrDefault() string {
	return configutil.CoalesceString(c.Addr, DefaultAddr)
}

// MustParseAddr returns the addr as a url.
func (c Config) MustParseAddr() *url.URL {
	remote, err := url.ParseRequestURI(c.AddrOrDefault())
	if err != nil {
		panic(err)
	}
	return remote
}

// MountOrDefault returns secrets mount or a default.
func (c Config) MountOrDefault() string {
	return configutil.CoalesceString(c.Mount, DefaultMount)
}

// TimeoutOrDefault returns the client timeout.
func (c Config) TimeoutOrDefault() time.Duration {
	return configutil.CoalesceDuration(c.Timeout, DefaultTimeout)
}
