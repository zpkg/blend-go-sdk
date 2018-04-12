package cron

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
)

func TestIsJobCanceled(t *testing.T) {
	assert := assert.New(t)

	ctx, cancel := context.WithCancel(context.Background())
	assert.False(IsJobCancelled(ctx))

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		assert.True(IsJobCancelled(ctx))
	}()
	cancel()
	wg.Wait()
}

func TestJobBuilder(t *testing.T) {
	assert := assert.New(t)

	assert.NotNil(NewJob().Schedule())
	assert.True(NewJob().IsEnabled())
	assert.True(NewJob().ShowMessages())
	assert.Equal("test_job", NewJob().WithName("test_job").Name())
	assert.Equal(EveryMinute(), NewJob().WithSchedule(EveryMinute()).Schedule())
	assert.Equal(time.Second, NewJob().WithTimeout(time.Second).Timeout())

	action := TaskAction(func(ctx context.Context) error { return nil })
	assert.NotNil(NewJob().WithAction(action).Action())

	assert.False(NewJob().WithIsEnabledProvider(func() bool { return false }).IsEnabled())
	assert.False(NewJob().WithShowMessagesProvider(func() bool { return false }).ShowMessages())
}
