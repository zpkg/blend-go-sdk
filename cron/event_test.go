package cron

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestNewEvent(t *testing.T) {
	assert := assert.New(t)

	e := NewEvent(FlagComplete, "test_task")
	assert.Equal(FlagComplete, e.Flag())
	assert.Equal("test_task", e.JobName)
	assert.False(e.Timestamp().IsZero())
	assert.True(e.IsEnabled())
	assert.True(e.IsWritable())
}
