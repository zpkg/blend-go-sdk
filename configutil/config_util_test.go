package configutil

import (
	"bytes"
	"testing"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/env"
)

type config struct {
	Environment string `json:"env" yaml:"env" env:"SERVICE_ENV"`
	Other       string `json:"other" yaml:"other" env:"OTHER"`
}

func TestRead(t *testing.T) {
	assert := assert.New(t)
	defer env.Restore()

	env.Env().Set(env.VarServiceEnv, "dev")
	var cfg config
	err := ReadFromReader(&cfg, bytes.NewBuffer([]byte("env: test\nother: foo")), ExtensionYAML)
	assert.Nil(err)
	assert.Equal("dev", cfg.Environment)
}

func TestReadPathUnset(t *testing.T) {
	assert := assert.New(t)
	defer env.Restore()

	env.Env().Set(env.VarServiceEnv, "dev")
	var cfg config
	err := ReadFromPath(&cfg, "")
	assert.True(IsPathUnset(err))
	assert.Equal("dev", cfg.Environment)
}
