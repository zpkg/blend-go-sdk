package cron

import (
	"context"
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/graceful"
	"github.com/blend/go-sdk/uuid"
)

var (
	_ graceful.Graceful = (*JobScheduler)(nil)
)

func TestJobSchedulerCullHistoryMaxAge(t *testing.T) {
	assert := assert.New(t)

	js := NewJobScheduler(&Config{
		History: HistoryConfig{
			MaxCount: 10,
			MaxAge:   6 * time.Hour,
		},
	}, NewJob("foo", noop))

	js.History = []JobInvocation{
		{ID: uuid.V4().String(), Started: time.Now().Add(-10 * time.Hour)},
		{ID: uuid.V4().String(), Started: time.Now().Add(-9 * time.Hour)},
		{ID: uuid.V4().String(), Started: time.Now().Add(-8 * time.Hour)},
		{ID: uuid.V4().String(), Started: time.Now().Add(-7 * time.Hour)},
		{ID: uuid.V4().String(), Started: time.Now().Add(-6 * time.Hour)},
		{ID: uuid.V4().String(), Started: time.Now().Add(-5 * time.Hour)},
		{ID: uuid.V4().String(), Started: time.Now().Add(-4 * time.Hour)},
		{ID: uuid.V4().String(), Started: time.Now().Add(-3 * time.Hour)},
		{ID: uuid.V4().String(), Started: time.Now().Add(-2 * time.Hour)},
		{ID: uuid.V4().String(), Started: time.Now().Add(-1 * time.Hour)},
	}

	filtered := js.cullHistory()
	assert.Len(filtered, 5)
}

func TestJobSchedulerCullHistoryMaxCount(t *testing.T) {
	assert := assert.New(t)

	js := NewJobScheduler(&Config{
		History: HistoryConfig{
			MaxCount: 5,
			MaxAge:   6 * time.Hour,
		},
	}, NewJob("foo", noop))

	js.History = []JobInvocation{
		{ID: uuid.V4().String(), Started: time.Now().Add(-10 * time.Minute)},
		{ID: uuid.V4().String(), Started: time.Now().Add(-9 * time.Minute)},
		{ID: uuid.V4().String(), Started: time.Now().Add(-8 * time.Minute)},
		{ID: uuid.V4().String(), Started: time.Now().Add(-7 * time.Minute)},
		{ID: uuid.V4().String(), Started: time.Now().Add(-6 * time.Minute)},
		{ID: uuid.V4().String(), Started: time.Now().Add(-5 * time.Minute)},
		{ID: uuid.V4().String(), Started: time.Now().Add(-4 * time.Minute)},
		{ID: uuid.V4().String(), Started: time.Now().Add(-3 * time.Minute)},
		{ID: uuid.V4().String(), Started: time.Now().Add(-2 * time.Minute)},
		{ID: uuid.V4().String(), Started: time.Now().Add(-1 * time.Minute)},
	}

	filtered := js.cullHistory()
	assert.Len(filtered, 5)
}

func TestJobSchedulerEnableDisable(t *testing.T) {
	assert := assert.New(t)

	var enabled, disabled bool
	js := NewJobScheduler(&Config{
		History: HistoryConfig{
			MaxCount: 5,
			MaxAge:   6 * time.Hour,
		},
	}, NewJob("foo", noop).WithOnDisabled(func(_ context.Context) {
		disabled = true
	}).WithOnEnabled(func(_ context.Context) {
		enabled = true
	}))

	js.Disable()
	assert.True(js.Disabled)

	assert.False(js.canRun())

	js.Enable()
	assert.False(js.Disabled)

	assert.True(disabled)
	assert.True(enabled)
}
