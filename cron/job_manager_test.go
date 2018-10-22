package cron

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/graceful"
	logger "github.com/blend/go-sdk/logger"
	"github.com/blend/go-sdk/uuid"
)

// assert the job manager is graceful
var (
	_ graceful.Graceful = (*JobManager)(nil)
)

func TestRunJobBySchedule(t *testing.T) {
	a := assert.New(t)

	didRun := make(chan struct{})

	jm := New()
	runAt := Now().Add(DefaultHeartbeatInterval)
	err := jm.LoadJob(&runAtJob{
		RunAt: runAt,
		RunDelegate: func(ctx context.Context) error {
			close(didRun)
			return nil
		},
	})
	a.Nil(err)

	jm.Start()
	defer jm.Stop()

	before := Now()
	<-didRun

	a.True(Since(before) < 2*DefaultHeartbeatInterval)
}

func TestDisableJob(t *testing.T) {
	a := assert.New(t)

	jm := New()
	a.Nil(jm.LoadJob(&runAtJob{RunAt: time.Now().UTC().Add(100 * time.Millisecond), RunDelegate: func(ctx context.Context) error {
		return nil
	}}))
	a.Nil(jm.DisableJob(runAtJobName))
	a.True(jm.IsDisabled(runAtJobName))
}

// The goal with this test is to see if panics take down the test process or not.
func TestJobManagerJobPanicHandling(t *testing.T) {
	assert := assert.New(t)

	manager := New()
	waitGroup := sync.WaitGroup{}
	waitGroup.Add(1)

	action := func(ctx context.Context) error {
		defer waitGroup.Done()
		panic("this is only a test")
	}
	manager.LoadJob(NewJob("panic-test").WithAction(action))
	manager.RunJob("panic-test")
	waitGroup.Wait()
	assert.True(true, "should complete")
}

func TestEnabledProvider(t *testing.T) {
	a := assert.New(t)
	a.StartTimeout(2000 * time.Millisecond)
	defer a.EndTimeout()

	manager := New()
	job := &testWithEnabled{
		isEnabled: true,
		action:    func() {},
	}

	manager.LoadJob(job)
	a.False(manager.IsDisabled("testWithEnabled"))
	manager.DisableJob("testWithEnabled")
	a.True(manager.IsDisabled("testWithEnabled"))
	job.isEnabled = false
	a.True(manager.IsDisabled("testWithEnabled"))
	manager.EnableJob("testWithEnabled")
	a.True(manager.IsDisabled("testWithEnabled"))
}

func TestFiresErrorOnTaskError(t *testing.T) {
	a := assert.New(t)
	a.StartTimeout(2000 * time.Millisecond)
	defer a.EndTimeout()

	agent := logger.New(logger.Error)
	defer agent.Close()

	manager := New().WithLogger(agent)
	defer manager.Stop()

	var errorDidFire bool
	var errorMatched bool
	wg := sync.WaitGroup{}
	wg.Add(2)

	agent.Listen(logger.Error, "foo", func(e logger.Event) {
		defer wg.Done()
		errorDidFire = true
		if typed, isTyped := e.(*logger.ErrorEvent); isTyped {
			if typed.Err() != nil {
				errorMatched = typed.Err().Error() == "this is only a test"
			}
		}
	})

	manager.LoadJob(NewJob("error_test").WithAction(func(ctx context.Context) error {
		defer wg.Done()
		return fmt.Errorf("this is only a test")
	}))
	manager.RunJob("error_test")
	wg.Wait()

	a.True(errorDidFire)
	a.True(errorMatched)
}

func TestManagerTracer(t *testing.T) {
	assert := assert.New(t)

	wg := sync.WaitGroup{}
	wg.Add(2)
	var didCallStart, didCallFinish bool
	var startTaskCorrect, finishTaskCorrect, errorUnset bool
	manager := New().
		WithTracer(&mockTracer{
			OnStart: func(ji *JobInvocation) {
				defer wg.Done()
				didCallStart = true
				startTaskCorrect = ji.Name == "tracer-test"
			},
			OnFinish: func(ji *JobInvocation) {
				defer wg.Done()
				didCallFinish = true
				finishTaskCorrect = ji.Name == "tracer-test"
				errorUnset = ji.Err == nil
			},
		})

	manager.LoadJob(NewJob("tracer-test"))
	manager.RunJob("tracer-test")
	wg.Wait()
	assert.True(didCallStart)
	assert.True(didCallFinish)
	assert.True(startTaskCorrect)
	assert.True(finishTaskCorrect)
	assert.True(errorUnset)
}

func TestJobManagerRunJobs(t *testing.T) {
	assert := assert.New(t)

	jm := New()
	jm.Start()
	defer jm.Stop()

	job0 := uuid.V4().String()
	job1 := uuid.V4().String()
	job2 := uuid.V4().String()

	var job0ran bool
	var job1ran bool
	var job2ran bool

	wg := sync.WaitGroup{}
	wg.Add(2)
	assert.Nil(jm.LoadJob(NewJob(job0).WithAction(func(_ context.Context) error {
		defer wg.Done()
		job0ran = true
		return nil
	})))
	assert.Nil(jm.LoadJob(NewJob(job1).WithAction(func(_ context.Context) error {
		defer wg.Done()
		job1ran = true
		return nil
	})))
	assert.Nil(jm.LoadJob(NewJob(job2).WithAction(func(_ context.Context) error {
		defer wg.Done()
		job2ran = true
		return nil
	})))

	assert.Nil(jm.RunJobs(job0, job2))
	wg.Wait()

	assert.True(job0ran)
	assert.False(job1ran)
	assert.True(job2ran)
}

func TestJobManagerRunAllJobs(t *testing.T) {
	assert := assert.New(t)

	jm := New()
	jm.Start()
	defer jm.Stop()

	job0 := uuid.V4().String()
	job1 := uuid.V4().String()
	job2 := uuid.V4().String()

	var job0ran bool
	var job1ran bool
	var job2ran bool

	wg := sync.WaitGroup{}
	wg.Add(3)
	assert.Nil(jm.LoadJob(NewJob(job0).WithAction(func(_ context.Context) error {
		defer wg.Done()
		job0ran = true
		return nil
	})))
	assert.Nil(jm.LoadJob(NewJob(job1).WithAction(func(_ context.Context) error {
		defer wg.Done()
		job1ran = true
		return nil
	})))
	assert.Nil(jm.LoadJob(NewJob(job2).WithAction(func(_ context.Context) error {
		defer wg.Done()
		job2ran = true
		return nil
	})))

	jm.RunAllJobs()
	wg.Wait()

	assert.True(job0ran)
	assert.True(job1ran)
	assert.True(job2ran)
}

func TestJobManagerJobLifecycle(t *testing.T) {
	assert := assert.New(t)

	jm := New()
	jm.Start()
	defer jm.Stop()

	var shouldFail bool
	j := newBrokenFixedTest(func(_ context.Context) error {
		defer func() {
			shouldFail = !shouldFail
		}()
		if shouldFail {
			return fmt.Errorf("only a test")
		}
		return nil
	})
	jm.LoadJob(j)

	assert.Nil(jm.RunJob("broken-fixed"))
	assert.Nil(jm.RunJob("broken-fixed"))
	<-j.BrokenSignal
	assert.Nil(jm.RunJob("broken-fixed"))
	<-j.FixedSignal

	assert.Equal(3, j.Starts)
	assert.Equal(1, j.Failures)
	assert.Equal(2, j.Completes)
}

func TestJobManagerJob(t *testing.T) {
	assert := assert.New(t)

	jm := New()
	j := newBrokenFixedTest(func(_ context.Context) error {
		return nil
	})
	jm.LoadJob(j)

	meta, err := jm.Job(j.Name())
	assert.Nil(err)
	assert.NotNil(meta)

	meta, err = jm.Job(uuid.V4().String())
	assert.NotNil(err)
	assert.Nil(meta)
}

func TestJobManagerLoadJob(t *testing.T) {
	assert := assert.New(t)

	jm := New()
	jm.LoadJob(&loadJobTestMinimum{})

	meta, err := jm.Job("load-job-test-minimum")
	assert.Nil(err)
	assert.NotNil(meta)

	assert.Equal("load-job-test-minimum", meta.Name)
	assert.NotNil(meta.Job)

	assert.NotNil(meta.EnabledProvider)
	assert.Equal(DefaultEnabled, meta.EnabledProvider())
	assert.NotNil(meta.SerialProvider)
	assert.Equal(DefaultSerial, meta.SerialProvider())
	assert.NotNil(meta.ShouldTriggerListenersProvider)
	assert.Equal(DefaultShouldTriggerListeners, meta.ShouldTriggerListenersProvider())
	assert.NotNil(meta.ShouldWriteOutputProvider)
	assert.Equal(DefaultShouldWriteOutput, meta.ShouldWriteOutputProvider())
}
