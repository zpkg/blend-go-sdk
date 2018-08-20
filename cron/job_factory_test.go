package cron

import (
	"context"
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
)

func TestJobFactory(t *testing.T) {
	assert := assert.New(t)

	assert.NotNil(NewJob("test_job").Schedule())
	assert.True(NewJob("test_job").IsEnabled())
	assert.True(NewJob("test_job").ShowMessages())
	assert.Equal("test_job", NewJob("test_job").Name())
	assert.Equal(EveryMinute(), NewJob("test_job").WithSchedule(EveryMinute()).Schedule())
	assert.Equal(time.Second, NewJob("test_job").WithTimeout(time.Second).Timeout())

	action := TaskAction(func(ctx context.Context) error { return nil })
	assert.NotNil(NewJob("test_job").WithAction(action).Action())

	assert.False(NewJob("test_job").WithIsEnabledProvider(func() bool { return false }).IsEnabled())
	assert.False(NewJob("test_job").WithShowMessagesProvider(func() bool { return false }).ShowMessages())
}
