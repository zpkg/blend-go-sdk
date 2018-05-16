package secrets

import (
	"net/url"
	"time"

	"github.com/blend/go-sdk/env"
	"github.com/blend/go-sdk/util"
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
	// Addr is the remote address of the secret store.
	Addr string `json:"addr" yaml:"addr" env:"VAULT_ADDR"`
	// Token is the authentication token used to talk to the secret store.
	Token string `json:"token" yaml:"token" env:"VAULT_TOKEN"`
	// Mount is the default mount path, it prefixes any keys.
	Mount string `json:"mount" yaml:"mount"`
	// Timeout is the dial timeout for requests to the secrets store.
	Timeout time.Duration `json:"timeout" yaml:"timeout"`
	// RootCAs is a list of certificate authority paths.
	RootCAs []string `json:"rootCAs" yaml:"rootCAs"`
}

// GetAddr returns the client addr.
func (c Config) GetAddr(inherited ...string) string {
	return util.Coalesce.String(c.Addr, DefaultAddr)
}

// MustRemote returns the addr as a url.
func (c Config) MustRemote() *url.URL {
	remote, err := url.ParseRequestURI(c.GetAddr())
	if err != nil {
		panic(err)
	}
	return remote
}

// GetToken returns the client token.
func (c Config) GetToken() string {
	return util.Coalesce.String(c.Token, "")
}

// GetMount returns the client token.
func (c Config) GetMount() string {
	return util.Coalesce.String(c.Mount, DefaultMount)
}

// GetTimeout returns the client timeout.
func (c Config) GetTimeout() time.Duration {
	return util.Coalesce.Duration(c.Timeout, DefaultTimeout)
}

// GetRootCAs returns root ca paths.
func (c Config) GetRootCAs() []string {
	return util.Coalesce.Strings(c.RootCAs, nil)
}
