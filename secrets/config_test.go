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

	assert.Equal("http://127.0.0.1:8100", cfg.MustRemote().String())
}
