package configutil

import (
	"github.com/blend/go-sdk/env"
)

// Env returns an environment variable source provider.
func Env(envVar string) ValueSource {
	return EnvValue(envVar)
}

// EnvValue returns a value from an environment variable.
type EnvValue string

// Value returns a given environment value.
func (ev EnvValue) Value() (string, error) {
	return env.Env().String(string(ev)), nil
}
