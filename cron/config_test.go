package cron

import (
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/env"
)

func TestNewConfigFromEnv(t *testing.T) {
	assert := assert.New(t)
	defer env.Restore()

	env.Env().Set(EnvVarHeartbeatInterval, "1s")

	cfg := NewConfigFromEnv()
	assert.Equal(time.Second, cfg.GetHeartbeatInterval())
}

func TestConfig(t *testing.T) {
	assert := assert.New(t)

	c := &Config{}

	assert.Equal(DefaultHeartbeatInterval, c.GetHeartbeatInterval())
	assert.Equal(DefaultHighPrecisionHeartbeatInterval, c.GetHeartbeatInterval(DefaultHighPrecisionHeartbeatInterval))

	set := &Config{HeartbeatInterval: time.Second}
	assert.Equal(time.Second, set.GetHeartbeatInterval(DefaultHighPrecisionHeartbeatInterval))
}
