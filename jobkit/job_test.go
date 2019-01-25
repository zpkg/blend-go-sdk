package jobkit

import (
	"context"
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/cron"
)

func TestJobProperties(t *testing.T) {
	assert := assert.New(t)

	job := NewJob(func(ctx context.Context) error {
		return nil
	})
	assert.NotNil(job.action)

	assert.NotEmpty(job.Name())
	job.WithName("foo")
	assert.Equal("foo", job.Name())

	assert.Nil(job.Schedule())
	job.WithSchedule(cron.EverySecond())
	assert.NotNil(job.Schedule())

	assert.Nil(job.NotificationsConfig())
	job.WithNotificationsConfig(&NotificationsConfig{})
	assert.NotNil(job.NotificationsConfig())

	assert.Zero(job.Timeout())
	job.WithTimeout(time.Second)
	assert.Equal(time.Second, job.Timeout())
}
