package configutil

import (
	"github.com/blend/go-sdk/env"
)

// EnvSource returns an environment variable source provider.
func EnvSource(envVar string) ValueSource {
	return EnvSourceValue(envVar)
}

// EnvSourceValue returns a value from an environment variable.
type EnvSourceValue string

// Value returns a given environment value.
func (esv EnvSourceValue) Value() (string, error) {
	return env.Env().String(string(esv)), nil
}
