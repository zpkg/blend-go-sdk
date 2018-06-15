package secrets

import (
	"testing"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/env"
)

func TestNewConfigFromEnv(t *testing.T) {
	assert := assert.New(t)
	defer env.Restore()

	env.Env().Set("VAULT_ADDR", "http://127.0.0.1:8100")
	env.Env().Set("VAULT_TOKEN", "thisisatest")

	cfg := NewConfigFromEnv()
	assert.Equal("http://127.0.0.1:8100", cfg.GetAddr())
	assert.Equal("thisisatest", cfg.GetToken())

	assert.Equal("http://127.0.0.1:8100", cfg.MustAddr().String())
}

func TestConfigIsZero(t *testing.T) {
	assert := assert.New(t)

	assert.True(Config{}.IsZero())
	assert.False(Config{Token: "garbage"}.IsZero())
}

func TestConfig(t *testing.T) {
	assert := assert.New(t)

	cfg := Config{}
	assert.Equal(DefaultAddr, cfg.GetAddr())
	assert.Empty(cfg.GetToken())
	assert.Equal(DefaultMount, cfg.GetMount())
	assert.Equal(DefaultTimeout, cfg.GetTimeout())
	assert.Empty(cfg.GetRootCAs())
}
