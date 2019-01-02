package cron

import (
	"context"
	"sync"
	"time"

	"github.com/blend/go-sdk/async"
	"github.com/blend/go-sdk/exception"
	"github.com/blend/go-sdk/logger"
)

// NewJobScheduler returns a job scheduler for a given job.
func NewJobScheduler(job Job) *JobScheduler {
	js := &JobScheduler{
		Latch: &async.Latch{},
		Name:  job.Name(),
		Job:   job,
	}

	if typed, ok := job.(ScheduleProvider); ok {
		js.Schedule = typed.Schedule()
	}

	if typed, ok := job.(TimeoutProvider); ok {
		js.TimeoutProvider = typed.Timeout
	} else {
		js.TimeoutProvider = func() time.Duration { return 0 }
	}

	if typed, ok := job.(EnabledProvider); ok {
		js.EnabledProvider = typed.Enabled
	} else {
		js.EnabledProvider = func() bool { return DefaultEnabled }
	}

	if typed, ok := job.(SerialProvider); ok {
		js.SerialProvider = typed.Serial
	} else {
		js.SerialProvider = func() bool { return DefaultSerial }
	}

	if typed, ok := job.(ShouldTriggerListenersProvider); ok {
		js.ShouldTriggerListenersProvider = typed.ShouldTriggerListeners
	} else {
		js.ShouldTriggerListenersProvider = func() bool { return DefaultShouldTriggerListeners }
	}

	if typed, ok := job.(ShouldWriteOutputProvider); ok {
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

	Tracer Tracer
	Log    *logger.Logger

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

// WithTracer sets the scheduler tracer.
func (js *JobScheduler) WithTracer(tracer Tracer) *JobScheduler {
	js.Tracer = tracer
	return js
}

// WithLogger sets the scheduler logger.
func (js *JobScheduler) WithLogger(log *logger.Logger) *JobScheduler {
	js.Log = log
	return js
}

// Start starts the scheduler.
func (js *JobScheduler) Start() {
	if !js.Latch.CanStart() {
		return
	}
	js.Latch.Starting()
	go js.runLoop()
	<-js.Latch.NotifyStarted()
}

// Stop stops the scheduler.
func (js *JobScheduler) Stop() {
	if !js.Latch.CanStop() {
		return
	}
	js.Latch.Stopping()
	<-js.Latch.NotifyStopped()
}

// Enable sets the job as enabled.
func (js *JobScheduler) Enable() {
	js.Lock()
	js.Disabled = false
	js.Unlock()
}

// Disable sets the job as disabled.
func (js *JobScheduler) Disable() {
	js.Lock()
	js.Disabled = true
	js.Unlock()
}

// Cancel stops an execution in process.
func (js *JobScheduler) Cancel() {
	if js.Current != nil {
		js.Current.Cancel()
	}
}

// RunLoop is the main scheduler loop.
// it alarms on the next runtime and forks a new routine to run the job.
// It can be aborted with the scheduler's async.Latch.
func (js *JobScheduler) RunLoop() {
	js.Latch.Started()

	// sniff the schedule, see if a next runtime is called for (or if the job is on demand).
	js.NextRuntime = Deref(js.Schedule.Next(Ref(js.NextRuntime)))
	if js.NextRuntime.IsZero() {
		js.Latch.Stopped()
		return
	}

	for {
		runAt := time.After(js.NextRuntime.UTC().Sub(Now()))
		select {
		case <-runAt:
			// start the job
			go js.Run()
			// set up the next runtime.
			js.NextRuntime = Deref(js.Schedule.Next(Ref(js.NextRuntime)))
		case <-js.Latch.NotifyStopping():
			js.Latch.Stopped()
			return
		}
	}
}

// Run forces the job to run.
// It checks if the job should be allowed to execute.
// It blocks on the job execution to enforce or clear timeouts.
func (js *JobScheduler) Run() {
	// check if the job is disabled.
	// check if the enabled provider is preventing the job from running.
	// check if the job is serial and currently executing.
	if !js.canRun() {
		return
	}

	// mark the start time
	start := Now()

	// create the root context.
	ctx, cancel := js.createContextWithCancel()

	// create a job invocation, or a record of each
	// individual execution of a job.
	ji := JobInvocation{
		ID:        NewJobInvocationID(),
		Name:      js.Name,
		StartTime: start,
		Context:   ctx,
		Cancel:    cancel,
	}

	js.Current = &ji
	defer func() {
		js.Current = nil
		js.Last = &ji
	}()

	finished := make(chan struct{})
	go js.execute(WithJobInvocation(ctx, &ji), &ji, finished)

	// check if we have to respect a timeout.
	if js.TimeoutProvider != nil {
		if timeout := js.TimeoutProvider(); timeout > 0 {
			ji.Timeout = start.Add(timeout)
			timeoutAlarm := time.After(timeout)
			for {
				select {
				case <-timeoutAlarm:
					cancel()
				case <-finished:
					return
				}
			}
		}
	}
	<-finished
}

//
// utility functions
//

// execute runs a given job invocation.
// it will signal lifecycle hooks.
func (js *JobScheduler) execute(ctx context.Context, ji *JobInvocation, finished chan struct{}) {
	var err error
	var tf TraceFinisher

	defer func() {
		if r := recover(); r != nil {
			err = exception.New(err)
		}
	}()

	if js.Tracer != nil {
		ctx, tf = js.Tracer.Start(ctx)
	}

	js.onStart(ctx, ji)

	// check if the job has been canceled.
	select {
	case <-ctx.Done():
		err = ErrJobCancelled
	case err = <-js.safeAsyncExec(ctx):
	}
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

	// signal to the waiters that the job is done executing.
	close(finished)
}

// safeAsyncExec runs a given job's body and recovers panics.
func (js *JobScheduler) safeAsyncExec(ctx context.Context) chan error {
	errors := make(chan error, 2)
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

func (js *JobScheduler) createContextWithCancel() (context.Context, context.CancelFunc) {
	return context.WithCancel(context.Background())
}

// canRun returns if a job can execute.
func (js *JobScheduler) canRun() bool {
	js.Lock()
	defer js.Unlock()

	if js.Disabled {
		return false
	}

	if js.EnabledProvider != nil {
		if !js.EnabledProvider() {
			return false
		}
	}

	if js.SerialProvider != nil && js.SerialProvider() {
		if js.Current != nil {
			return false
		}
	}
	return true
}

func (js *JobScheduler) onStart(ctx context.Context, ji *JobInvocation) {
	if js.Log != nil && js.ShouldTriggerListenersProvider() {
		event := NewEvent(FlagStarted, ji.Name).WithIsWritable(js.ShouldWriteOutputProvider())
		js.Log.SubContext(ji.ID).Trigger(event)
	}
	if typed, ok := js.Job.(OnStartReceiver); ok {
		typed.OnStart(ctx)
	}
}

func (js *JobScheduler) onCancelled(ctx context.Context, ji *JobInvocation) {
	if js.Log != nil && js.ShouldTriggerListenersProvider() {
		event := NewEvent(FlagCancelled, ji.Name).
			WithIsWritable(js.ShouldWriteOutputProvider()).
			WithElapsed(ji.Elapsed)
		js.Log.SubContext(ji.ID).Trigger(event)
	}
	if typed, ok := js.Job.(OnCancellationReceiver); ok {
		typed.OnCancellation(ctx)
	}
}

func (js *JobScheduler) onComplete(ctx context.Context, ji *JobInvocation) {
	if js.Log != nil && js.ShouldTriggerListenersProvider() {
		event := NewEvent(FlagComplete, ji.Name).
			WithIsWritable(js.ShouldWriteOutputProvider()).
			WithElapsed(ji.Elapsed)
		js.Log.SubContext(ji.ID).Trigger(event)
	}
	if typed, ok := js.Job.(OnCompleteReceiver); ok {
		typed.OnComplete(ctx)
	}

	if js.Last != nil && js.Last.Err != nil {
		if js.Log != nil {
			event := NewEvent(FlagFixed, ji.Name).
				WithIsWritable(js.ShouldWriteOutputProvider()).
				WithElapsed(ji.Elapsed)
			js.Log.SubContext(ji.ID).Trigger(event)
		}

		if typed, ok := js.Job.(OnFixedReceiver); ok {
			typed.OnFixed(ctx)
		}
	}
}

func (js *JobScheduler) onFailure(ctx context.Context, ji *JobInvocation) {
	if js.Log != nil && js.ShouldTriggerListenersProvider() {
		event := NewEvent(FlagFailed, ji.Name).
			WithIsWritable(js.ShouldWriteOutputProvider()).
			WithElapsed(ji.Elapsed).
			WithErr(ji.Err)

		js.Log.SubContext(ji.ID).Trigger(event)
	}
	if ji.Err != nil {
		logger.MaybeError(js.Log, ji.Err)
	}
	if typed, ok := js.Job.(OnFailureReceiver); ok {
		typed.OnFailure(ctx)
	}
	if js.Last != nil && js.Last.Err == nil {
		if js.Log != nil {
			event := NewEvent(FlagBroken, ji.Name).
				WithIsWritable(js.ShouldWriteOutputProvider()).
				WithElapsed(ji.Elapsed)
			js.Log.SubContext(ji.ID).Trigger(event)
		}

		if typed, ok := js.Job.(OnBrokenReceiver); ok {
			typed.OnBroken(ctx)
		}
	}
}
