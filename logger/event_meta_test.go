package logger

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestNewEventMeta(t *testing.T) {
	assert := assert.New(t)

	em := NewEventMeta(Info)
	assert.Equal(Info, em.Flag())
	assert.False(em.Timestamp().IsZero())
}
