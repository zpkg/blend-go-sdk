package cron

import (
	"context"
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
)

func TestJobBuilder(t *testing.T) {
	assert := assert.New(t)

	assert.NotNil(NewJob("test_job"))
	assert.True(NewJob("test_job").Enabled())
	assert.Zero(NewJob("test_job").Timeout())
	assert.True(NewJob("test_job").ShouldTriggerListeners())
	assert.True(NewJob("test_job").ShouldWriteOutput())
	assert.Equal("test_job", NewJob("test_job").Name())
	assert.Equal("test_job2", NewJob("test_job").WithName("test_job2").Name())
	assert.Equal(EveryMinute(), NewJob("test_job").WithSchedule(EveryMinute()).Schedule())
	assert.Equal(time.Second, NewJob("test_job").WithTimeoutProvider(func() time.Duration { return time.Second }).Timeout())
	action := Action(func(ctx context.Context) error { return nil })
	assert.NotNil(NewJob("test_job").WithAction(action).action)

	assert.False(NewJob("test_job").WithEnabledProvider(func() bool { return false }).Enabled())
	assert.False(NewJob("test_job").WithShouldTriggerListenersProvider(func() bool { return false }).ShouldTriggerListeners())
	assert.False(NewJob("test_job").WithShouldWriteOutputProvider(func() bool { return false }).ShouldWriteOutput())
}

func TestJobBuilderLifecycle(t *testing.T) {
	assert := assert.New(t)
	job := NewJob("test_job")

	assert.Nil(job.onStart)
	var onStart bool
	job.WithOnStart(func(ji *JobInvocation) {
		onStart = true
	})
	assert.NotNil(job.onStart)
	job.OnStart(nil) // this will break context handling code if not nil checked.
	assert.True(onStart)

	assert.Nil(job.onCancellation)
	var onCancellation bool
	job.WithOnCancellation(func(ji *JobInvocation) {
		onCancellation = true
	})
	assert.NotNil(job.onCancellation)
	job.OnCancellation(nil)
	assert.True(onCancellation)

	assert.Nil(job.onComplete)
	var onComplete bool
	job.WithOnComplete(func(ji *JobInvocation) {
		onComplete = true
	})
	assert.NotNil(job.onComplete)
	job.OnComplete(nil)
	assert.True(onComplete)

	assert.Nil(job.onFailure)
	var onFailure bool
	job.WithOnFailure(func(ji *JobInvocation) {
		onFailure = true
	})
	assert.NotNil(job.onFailure)
	job.OnFailure(nil)
	assert.True(onFailure)

	assert.Nil(job.onBroken)
	var onBroken bool
	job.WithOnBroken(func(ji *JobInvocation) {
		onBroken = true
	})
	assert.NotNil(job.onBroken)
	job.OnBroken(nil)
	assert.True(onBroken)

	assert.Nil(job.onFixed)
	var onFixed bool
	job.WithOnFixed(func(ji *JobInvocation) {
		onFixed = true
	})
	assert.NotNil(job.onFixed)
	job.OnFixed(nil)
	assert.True(onFixed)
}
