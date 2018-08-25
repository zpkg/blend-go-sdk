package stats

import (
	"testing"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/logger"
)

func TestAddWebListeners(t *testing.T) {
	assert := assert.New(t)

	log := logger.None()
	AddWebListeners(nil, nil)
	assert.False(log.HasListener(logger.HTTPResponse, "stats"))
	AddWebListeners(log, NewMockCollector())
	assert.True(log.HasListener(logger.HTTPResponse, "stats"))
}
