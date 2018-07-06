package cron

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/exception"
	logger "github.com/blend/go-sdk/logger"
)

const (
	runAtJobName = "runAt"
)

type runAtJob struct {
	RunAt       time.Time
	RunDelegate func(ctx context.Context) error
}

type runAt time.Time

func (ra runAt) GetNextRunTime(after *time.Time) *time.Time {
	typed := time.Time(ra)
	return &typed
}

func (raj *runAtJob) Name() string {
	return "runAt"
}

func (raj *runAtJob) Schedule() Schedule {
	return runAt(raj.RunAt)
}

func (raj *runAtJob) Execute(ctx context.Context) error {
	return raj.RunDelegate(ctx)
}

type testJobWithTimeout struct {
	RunAt                time.Time
	TimeoutDuration      time.Duration
	RunDelegate          func(ctx context.Context) error
	CancellationDelegate func()
}

func (tj *testJobWithTimeout) Name() string {
	return "testJobWithTimeout"
}

func (tj *testJobWithTimeout) Timeout() time.Duration {
	return tj.TimeoutDuration
}

func (tj *testJobWithTimeout) Schedule() Schedule {
	return Immediately()
}

func (tj *testJobWithTimeout) Execute(ctx context.Context) error {
	return tj.RunDelegate(ctx)
}

func (tj *testJobWithTimeout) OnCancellation() {
	tj.CancellationDelegate()
}

type testJobInterval struct {
	RunEvery    time.Duration
	RunDelegate func(ctx context.Context) error
}

func (tj *testJobInterval) Name() string {
	return "testJobInterval"
}

func (tj *testJobInterval) Schedule() Schedule {
	return Every(tj.RunEvery)
}

func (tj *testJobInterval) Execute(ctx context.Context) error {
	return tj.RunDelegate(ctx)
}

func TestRunTask(t *testing.T) {
	a := assert.New(t)
	a.StartTimeout(2000 * time.Millisecond)
	defer a.EndTimeout()

	didRun := make(chan struct{})
	New().RunTask(NewTask(func(ctx context.Context) error {
		close(didRun)
		return nil
	}))

	<-didRun
}

func TestRunTaskAndCancel(t *testing.T) {
	a := assert.New(t)

	jm := New()

	didRun := make(chan struct{})
	didFinish := make(chan struct{})
	jm.RunTask(NewTaskWithName("taskToCancel", func(ctx context.Context) error {
		defer func() {
			close(didFinish)
		}()
		close(didRun)
		alarm := time.After(time.Second)
		select {
		case <-ctx.Done():
			return nil
		case <-alarm:
			return exception.New("timed out")
		}
	}))

	<-didRun
	a.Nil(jm.CancelTask("taskToCancel"))
	<-didFinish
}

func TestRunJobBySchedule(t *testing.T) {
	a := assert.New(t)

	didRun := make(chan struct{})

	jm := New().WithHighPrecisionHeartbeat()
	runAt := Now().Add(jm.HeartbeatInterval())
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

	a.True(Since(before) < 2*jm.HeartbeatInterval())
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

func TestSerialTask(t *testing.T) {
	assert := assert.New(t)

	// test that serial execution actually blocks
	runCount := new(AtomicCounter)
	jm := New()

	task := NewSerialTaskWithName("test", func(ctx context.Context) error {
		runCount.Increment()
		time.Sleep(10 * time.Millisecond)
		return nil
	})
	jm.RunTask(task)
	jm.RunTask(task)
	time.Sleep(100 * time.Millisecond)
	assert.Equal(1, runCount.Get())

	// ensure parallel execution is still working as intended
	task = NewTaskWithName("test1", func(ctx context.Context) error {
		runCount.Increment()
		time.Sleep(10 * time.Millisecond)
		return nil
	})
	runCount = new(AtomicCounter)
	jm.RunTask(task)
	jm.RunTask(task)
	time.Sleep(100 * time.Millisecond)
	assert.Equal(2, runCount.Get())
}

func TestRunJobByScheduleRapid(t *testing.T) {
	a := assert.New(t)

	didRun := make(chan struct{})
	// high precision heartbeat is somewhere around 5ms
	jm := New().WithHighPrecisionHeartbeat()
	err := jm.LoadJob(&testJobInterval{RunEvery: time.Millisecond, RunDelegate: func(ctx context.Context) error {
		close(didRun)
		return nil
	}})
	a.Nil(err)

	jm.Start()
	defer jm.Stop()

	alarm := time.After(50 * time.Millisecond)
	select {
	case <-didRun:
		break
	case <-alarm:
		a.FailNow("timed out")
	}
}

// The goal with this test is to see if panics take down the test process or not.
func TestJobManagerTaskPanicHandling(t *testing.T) {
	a := assert.New(t)

	manager := New()
	waitGroup := sync.WaitGroup{}
	waitGroup.Add(1)
	err := manager.RunTask(NewTask(func(ctx context.Context) error {
		defer waitGroup.Done()
		panic("this is only a test")
	}))

	waitGroup.Wait()
	a.Nil(err)
}

type testWithEnabled struct {
	isEnabled bool
	action    func()
}

func (twe testWithEnabled) Name() string {
	return "testWithEnabled"
}

func (twe testWithEnabled) Schedule() Schedule {
	return OnDemand()
}

func (twe testWithEnabled) Enabled() bool {
	return twe.isEnabled
}

func (twe testWithEnabled) Execute(ctx context.Context) error {
	twe.action()
	return nil
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
	manager := New()
	manager.SetLogger(agent)

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

	manager.LoadJob(NewJob().WithAction(func(ctx context.Context) error {
		defer wg.Done()
		return fmt.Errorf("this is only a test")
	}).WithName("error_test"))
	manager.RunJob("error_test")
	wg.Wait()

	a.True(errorDidFire)
	a.True(errorMatched)
}
