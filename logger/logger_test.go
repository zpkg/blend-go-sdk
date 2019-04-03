package logger

import (
	"testing"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/uuid"
)

func TestNew(t *testing.T) {
	assert := assert.New(t)

	log, err := New()
	assert.Nil(err)
	assert.NotNil(log.Latch)
	assert.NotNil(log.Context)
	assert.NotNil(log.Formatter)
	assert.NotNil(log.Output)
	assert.True(log.RecoverPanics)

	for _, defaultFlag := range DefaultFlags {
		assert.True(log.Flags.IsEnabled(defaultFlag))
	}

	log, err = New(OptAll(), OptFormatter(NewJSONOutputFormatter()))
	assert.Nil(err)
	assert.True(log.Flags.IsEnabled(uuid.V4().String()))
	typed, ok := log.Formatter.(*JSONOutputFormatter)
	assert.True(ok)
	assert.NotNil(typed)
}
