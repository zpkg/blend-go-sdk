package cron

import (
	"context"
	"sync"
	"time"

	"github.com/blend/go-sdk/async"
	"github.com/blend/go-sdk/exception"
)

// NewJobScheduler returns a job scheduler for a given job.
func NewJobScheduler(job Job) *JobScheduler {
	js := JobScheduler{
		Latch: &async.Latch{},
		Name:  job.Name(),
	}

	if typed, ok := j.(ScheduleProvider); ok {
		js.Schedule = typed.Schedule()
	}

	if typed, ok := j.(TimeoutProvider); ok {
		js.TimeoutProvider = typed.Timeout
	} else {
		js.TimeoutProvider = func() time.Duration { return 0 }
	}

	if typed, ok := j.(EnabledProvider); ok {
		js.EnabledProvider = typed.Enabled
	} else {
		js.EnabledProvider = func() bool { return DefaultEnabled }
	}

	if typed, ok := j.(SerialProvider); ok {
		js.SerialProvider = typed.Serial
	} else {
		js.SerialProvider = func() bool { return DefaultSerial }
	}

	if typed, ok := j.(ShouldTriggerListenersProvider); ok {
		js.ShouldTriggerListenersProvider = typed.ShouldTriggerListeners
	} else {
		js.ShouldTriggerListenersProvider = func() bool { return DefaultShouldTriggerListeners }
	}

	if typed, ok := j.(ShouldWriteOutputProvider); ok {
		js.ShouldWriteOutputProvider = typed.ShouldWriteOutput
	} else {
		js.ShouldWriteOutputProvider = func() bool { return DefaultShouldWriteOutput }
	}

	return js
}

// JobScheduler is a job instance.
type JobScheduler struct {
	sync.Mutex

	Latch *async.Latch

	Name string
	Job  Job

	// Meta Fields
	Disabled    bool
	NextRuntime time.Time
	Current     *JobInvocation
	Last        *JobInvocation

	Schedule                       Schedule
	EnabledProvider                func() bool
	SerialProvider                 func() bool
	TimeoutProvider                func() time.Duration
	ShouldTriggerListenersProvider func() bool
	ShouldWriteOutputProvider      func() bool
}

func (js *JobScheduler) Start() {
	js.Latch.Starting()
	go func() {
		for {

		}
	}()
	<-js.Latch.NotifyStarted()
}

func (js *JobScheduler) Stop() {
	js.Latch.Stopping()
	<-js.Latch.NotifyStopped()
}

func (js *JobScheduler) Enable() {
	js.Lock()
	defer js.Unlock()
	js.Start()
	js.Disabled = false
}

func (js *JobScheduler) Disable() {
	js.Lock()
	defer js.Unlock()
	js.Stop()
	js.Disabled = true
}

func (js *JobScheduler) run() {
	if !js.canRun() {
		return
	}

	now := Now()
	js.NextRunTime = js.scheduleNextRuntime(js.Schedule, Optional(now))

	start := Now()
	ctx, cancel := jm.createContextWithCancel()

	ji := JobInvocation{
		ID:        NewJobInvocationID(),
		Name:      jobMeta.Name,
		StartTime: start,
		JobMeta:   jobMeta,
		Context:   ctx,
		Cancel:    cancel,
	}

	if jobMeta.TimeoutProvider != nil {
		if timeout := jobMeta.TimeoutProvider(); timeout > 0 {
			ji.Timeout = start.Add(timeout)
		}
	}
	go js.execute(WithJobInvocation(ctx, &ji), &ji)
}

func (js *JobScheduler) execute(ctx context.Context, ji *JobInvocation) {
	var err error
	var tf TraceFinisher
	defer func() {
		if tf != nil {
			tf.Finish(ctx)
		}
		ji.Elapsed = Since(ji.StartTime)
		ji.Err = err
		if err != nil && IsJobCancelled(err) {
			js.onCancelled(ctx, ji)
		} else if ji.Err != nil {
			js.onFailure(ctx, ji)
		} else {
			js.onComplete(ctx, ji)
		}
		js.JobMeta.Last = ji
	}()
	if js.tracer != nil {
		ctx, tf = js.tracer.Start(ctx)
	}

	js.onStart(ctx, ji)

	select {
	case <-ctx.Done():
		err = ErrJobCancelled
	case err = <-js.safeAsyncExec(ctx, ji.JobMeta.Job):
		return
	}
}

func (js *JobScheduler) safeAsyncExec(ctx context.Context) chan error {
	errors := make(chan error)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				errors <- exception.New(r)
			}
		}()
		errors <- js.Job.Execute(ctx)
	}()
	return errors
}

func (js *JobScheduler) canRun() bool {
	if js.Meta.Disabled {
		return false
	}
	if js.EnabledProvider != nil {
		if !js.EnabledProvider() {
			return false
		}
	}

	if js.SerialProvider != nil && js.SerialProvider() {
		if js.Executing != nil {
			return false
		}
	}
	return true
}
