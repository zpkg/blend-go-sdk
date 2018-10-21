package cron

import (
	"context"
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
)

func TestJobFactory(t *testing.T) {
	assert := assert.New(t)

	assert.NotNil(NewJob("test_job"))
	assert.True(NewJob("test_job").Enabled())
	assert.Zero(NewJob("test_job").Timeout())
	assert.True(NewJob("test_job").ShouldTriggerListeners())
	assert.True(NewJob("test_job").ShouldWriteOutput())
	assert.Equal("test_job", NewJob("test_job").Name())
	assert.Equal(EveryMinute(), NewJob("test_job").WithSchedule(EveryMinute()).Schedule())
	assert.Equal(time.Second, NewJob("test_job").WithTimeoutProvider(func() time.Duration { return time.Second }).Timeout())

	action := Action(func(ctx context.Context) error { return nil })
	assert.NotNil(NewJob("test_job").WithAction(action).Action())

	assert.False(NewJob("test_job").WithEnabledProvider(func() bool { return false }).Enabled())
	assert.False(NewJob("test_job").WithShouldTriggerListenersProvider(func() bool { return false }).ShouldTriggerListeners())
	assert.False(NewJob("test_job").WithShouldWriteOutputProvider(func() bool { return false }).ShouldWriteOutput())
}
