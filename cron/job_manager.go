package cron

// NOTE: ALL TIMES ARE IN UTC. JUST USE UTC.

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/blend/go-sdk/async"
	"github.com/blend/go-sdk/exception"
	"github.com/blend/go-sdk/logger"
)

// New returns a new job manager.
func New() *JobManager {
	jm := JobManager{
		latch:             &async.Latch{},
		heartbeatInterval: DefaultHeartbeatInterval,
		jobs:              map[string]*JobMeta{},
		tasks:             map[string]*TaskMeta{},
	}
	jm.schedulerWorker = async.NewInterval(jm.runDueJobs, DefaultHeartbeatInterval)
	jm.killHangingTasksWorker = async.NewInterval(jm.killHangingTasks, DefaultHeartbeatInterval)
	return &jm
}

// NewFromConfig returns a new job manager from a given config.
func NewFromConfig(cfg *Config) *JobManager {
	return New().WithHeartbeatInterval(cfg.GetHeartbeatInterval())
}

// NewFromEnv returns a new job manager from the environment.
func NewFromEnv() (*JobManager, error) {
	cfg, err := NewConfigFromEnv()
	if err != nil {
		return nil, err
	}
	return NewFromConfig(cfg), nil
}

// MustNewFromEnv returns a new job manager from the environment.
func MustNewFromEnv() *JobManager {
	cfg, err := NewConfigFromEnv()
	if err != nil {
		panic(err)
	}
	return NewFromConfig(cfg)
}

// JobManager is the main orchestration and job management object.
type JobManager struct {
	sync.Mutex

	latch  *async.Latch
	tracer Tracer

	heartbeatInterval time.Duration
	log               *logger.Logger

	schedulerWorker        *async.Interval
	killHangingTasksWorker *async.Interval

	jobs  map[string]*JobMeta
	tasks map[string]*TaskMeta
}

// Logger returns the diagnostics agent.
func (jm *JobManager) Logger() *logger.Logger {
	return jm.log
}

// WithLogger sets the logger and returns a reference to the job manager.
func (jm *JobManager) WithLogger(log *logger.Logger) *JobManager {
	jm.Lock()
	defer jm.Unlock()

	jm.log = log
	return jm
}

// WithTracer sets the manager's tracer.
func (jm *JobManager) WithTracer(tracer Tracer) *JobManager {
	jm.tracer = tracer
	return jm
}

// Tracer returns the manager's tracer.
func (jm *JobManager) Tracer() Tracer {
	return jm.tracer
}

// WithHighPrecisionHeartbeat sets the heartbeat interval to the high precision interval and returns the job manager.
func (jm *JobManager) WithHighPrecisionHeartbeat() *JobManager {
	return jm.WithHeartbeatInterval(DefaultHighPrecisionHeartbeatInterval)
}

// WithDefaultHeartbeat sets the heartbeat interval to the default interval and returns the job manager.
func (jm *JobManager) WithDefaultHeartbeat() *JobManager {
	return jm.WithHeartbeatInterval(DefaultHeartbeatInterval)
}

// WithHeartbeatInterval sets the heartbeat interval explicitly and returns the job manager.
func (jm *JobManager) WithHeartbeatInterval(interval time.Duration) *JobManager {
	jm.schedulerWorker.WithInterval(interval)
	jm.killHangingTasksWorker.WithInterval(interval)
	jm.heartbeatInterval = interval
	return jm
}

// HeartbeatInterval returns the current heartbeat interval.
func (jm *JobManager) HeartbeatInterval() time.Duration {
	return jm.heartbeatInterval
}

// ----------------------------------------------------------------------------
// Informational Methods
// ----------------------------------------------------------------------------

// HasJob returns if a jobName is loaded or not.
func (jm *JobManager) HasJob(jobName string) (hasJob bool) {
	jm.Lock()
	defer jm.Unlock()

	_, hasJob = jm.jobs[jobName]
	return
}

// Job returns a job instance by name.
func (jm *JobManager) Job(jobName string) (job Job) {
	jm.Lock()
	defer jm.Unlock()

	if jobMeta, hasJob := jm.jobs[jobName]; hasJob {
		job = jobMeta.Job
	}
	return
}

// IsDisabled returns if a job is disabled.
func (jm *JobManager) IsDisabled(jobName string) (value bool) {
	jm.Lock()
	defer jm.Unlock()

	if job, hasJob := jm.jobs[jobName]; hasJob {
		value = job.Disabled
		if job.EnabledProvider != nil {
			value = value || !job.EnabledProvider()
		}
	}
	return
}

// IsRunning returns if a task is currently running.
func (jm *JobManager) IsRunning(taskName string) (isRunning bool) {
	jm.Lock()
	defer jm.Unlock()

	_, isRunning = jm.tasks[taskName]
	return
}

// ReadAllJobs allows the consumer to do something with the full job list, using a read lock.
func (jm *JobManager) ReadAllJobs(action func(jobs map[string]*JobMeta)) {
	jm.Lock()
	defer jm.Unlock()

	action(jm.jobs)
}

// Status returns a status object.
func (jm *JobManager) Status() *Status {
	jm.Lock()
	defer jm.Unlock()

	status := Status{
		Tasks: map[string]TaskMeta{},
	}
	for _, meta := range jm.jobs {
		status.Jobs = append(status.Jobs, *meta)
	}
	for name, task := range jm.tasks {
		status.Tasks[name] = *task
	}
	return &status
}

// --------------------------------------------------------------------------------
// Core Methods
// --------------------------------------------------------------------------------

// LoadJobs loads a variadic list of jobs.
func (jm *JobManager) LoadJobs(jobs ...Job) error {
	jm.Lock()
	defer jm.Unlock()

	var err error
	for _, job := range jobs {
		err = jm.loadJobUnsafe(job)
		if err != nil {
			return err
		}
	}
	return nil
}

// LoadJob loads a job.
func (jm *JobManager) LoadJob(job Job) error {
	jm.Lock()
	defer jm.Unlock()

	return jm.loadJobUnsafe(job)
}

// DisableJobs disables a variadic list of job names.
func (jm *JobManager) DisableJobs(jobNames ...string) error {
	jm.Lock()
	defer jm.Unlock()

	var err error
	for _, jobName := range jobNames {
		err = jm.setJobDisabledUnsafe(jobName, true)
		if err != nil {
			return err
		}
	}
	return nil
}

// DisableJob stops a job from running but does not unload it.
func (jm *JobManager) DisableJob(jobName string) error {
	jm.Lock()
	defer jm.Unlock()

	return jm.setJobDisabledUnsafe(jobName, true)
}

// EnableJobs enables a variadic list of job names.
func (jm *JobManager) EnableJobs(jobNames ...string) error {
	var err error
	for _, jobName := range jobNames {
		err = jm.setJobDisabledUnsafe(jobName, false)
		if err != nil {
			return err
		}
	}
	return nil
}

// EnableJob enables a job that has been disabled.
func (jm *JobManager) EnableJob(jobName string) error {
	jm.Lock()
	defer jm.Unlock()

	return jm.setJobDisabledUnsafe(jobName, false)
}

// RunJobs runs a variadic list of job names.
func (jm *JobManager) RunJobs(jobNames ...string) error {
	jm.Lock()
	defer jm.Unlock()
	for _, jobName := range jobNames {
		if jobMeta, ok := jm.jobs[jobName]; ok {
			jm.runJobUnsafe(jobMeta)
		} else {
			return exception.New(ErrJobNotLoaded).WithMessagef("job: %s", jobName)
		}
	}
	return nil
}

// RunJob runs a job by jobName on demand.
func (jm *JobManager) RunJob(jobName string) error {
	jm.Lock()
	defer jm.Unlock()

	if job, ok := jm.jobs[jobName]; ok {
		jm.runJobUnsafe(job)
		return nil
	}
	return exception.New(ErrJobNotLoaded).WithMessagef("job: %s", jobName)
}

// RunAllJobs runs every job that has been loaded in the JobManager at once.
func (jm *JobManager) RunAllJobs() {
	jm.Lock()
	defer jm.Unlock()
	for _, jobMeta := range jm.jobs {
		jm.runJobUnsafe(jobMeta)
	}
}

// RunTask runs a task.
func (jm *JobManager) RunTask(task Task) {
	jm.Lock()
	defer jm.Unlock()
	jm.runTaskUnsafe(task)
}

// CancelTask cancels (sends the cancellation signal) to a running task.
func (jm *JobManager) CancelTask(taskName string) (err error) {
	jm.Lock()
	defer jm.Unlock()

	task, hasTask := jm.tasks[taskName]
	if !hasTask {
		err = exception.New(ErrTaskNotFound).WithMessagef("task: %s", taskName)
		return
	}
	jm.onTaskCancelled(task.Task, Since(task.StartTime))
	task.Cancel()
	return
}

// Start begins the schedule runner for a JobManager.
func (jm *JobManager) Start() error {
	if !jm.latch.CanStart() {
		return fmt.Errorf("already started")
	}
	jm.latch.Starting()
	jm.schedulerWorker.Start()
	jm.killHangingTasksWorker.Start()
	jm.latch.Started()
	return nil
}

// Stop stops the schedule runner for a JobManager.
func (jm *JobManager) Stop() error {
	if !jm.latch.CanStop() {
		return fmt.Errorf("already stopped")
	}
	jm.latch.Stopping()
	jm.schedulerWorker.Stop()
	jm.killHangingTasksWorker.Stop()
	jm.latch.Stopped()
	return nil
}

// NotifyStarted returns the started notification channel.
func (jm *JobManager) NotifyStarted() <-chan struct{} {
	return jm.latch.NotifyStarted()
}

// NotifyStopped returns the stopped notification channel.
func (jm *JobManager) NotifyStopped() <-chan struct{} {
	return jm.latch.NotifyStopped()
}

// --------------------------------------------------------------------------------
// lifecycle methods
// --------------------------------------------------------------------------------

func (jm *JobManager) runDueJobs() error {
	jm.Lock()
	defer jm.Unlock()
	now := Now()
	for _, jobMeta := range jm.jobs {
		if !jobMeta.NextRunTime.IsZero() && jobMeta.NextRunTime.Before(now) {
			jm.runJobUnsafe(jobMeta)
		}
	}
	return nil
}

func (jm *JobManager) killHangingTasks() (err error) {
	jm.Lock()
	defer jm.Unlock()

	var effectiveTimeout time.Time
	var now time.Time
	var t1, t2 time.Time

	for taskName, taskMeta := range jm.tasks {
		if taskMeta.Timeout.IsZero() {
			return
		}
		now = Now()
		if jobMeta, hasJobMeta := jm.jobs[taskName]; hasJobMeta {
			nextRuntime := jobMeta.NextRunTime
			t1 = taskMeta.Timeout
			t2 = nextRuntime
			effectiveTimeout = Min(t1, t2)
		} else {
			effectiveTimeout = taskMeta.Timeout
		}
		if effectiveTimeout.Before(now) {
			jm.killHangingTask(taskMeta)
		}
	}
	return nil
}

//
// these assume a lock is held
//

func (jm *JobManager) runJobUnsafe(jobMeta *JobMeta) {
	now := Now()

	// merge the runtime disabled with the job provided enabled/disabled
	disabled := jobMeta.Disabled
	if jobMeta.EnabledProvider != nil {
		disabled = disabled || !jobMeta.EnabledProvider()
	}

	if !disabled {
		jobMeta.NextRunTime = jm.scheduleNextRuntime(jobMeta.Schedule, Optional(now))
		jobMeta.LastRunTime = now
		jm.runTaskUnsafe(jobMeta.Job)
	}
}

func (jm *JobManager) runTaskUnsafe(t Task) {
	if !jm.shouldRunTask(t) {
		return
	}

	taskName := t.Name()
	start := Now()
	ctx, cancel := jm.createContext()

	tm := &TaskMeta{
		Name:      taskName,
		StartTime: start,
		Task:      t,
		Context:   ctx,
		Cancel:    cancel,
	}

	if typed, ok := t.(TaskTimeoutProvider); ok {
		tm.Timeout = start.Add(typed.Timeout())
	}
	jm.tasks[taskName] = tm

	// run the actual task
	go func() {
		var err error
		var tf TraceFinisher
		defer func() {
			if r := recover(); r != nil {
				err = exception.New(r)
			}
			if tf != nil {
				tf.Finish(ctx, t, err)
			}

			jm.Lock()
			if _, hasTask := jm.tasks[taskName]; hasTask {
				delete(jm.tasks, taskName)
			}
			jm.Unlock()

			jm.onTaskComplete(t, Since(start), err)
		}()
		if jm.tracer != nil {
			ctx, tf = jm.tracer.Start(ctx, t)
		}

		// fire the on start listener if any
		jm.onTaskStart(t)
		// run the task
		err = t.Execute(ctx)
	}()

	return
}

func (jm *JobManager) killHangingTask(task *TaskMeta) {
	task.Cancel()
	jm.onTaskCancelled(task.Task, Since(task.StartTime))
	delete(jm.tasks, task.Name)
}

// --------------------------------------------------------------------------------
// Utility Methods
// --------------------------------------------------------------------------------

// LoadJob adds a job to the manager.
func (jm *JobManager) loadJobUnsafe(j Job) error {
	jobName := j.Name()

	if _, hasJob := jm.jobs[jobName]; hasJob {
		return exception.New(ErrJobAlreadyLoaded).WithMessagef("job: %s", j.Name())
	}

	schedule := j.Schedule()
	meta := &JobMeta{
		Name:        jobName,
		Job:         j,
		NextRunTime: jm.scheduleNextRuntime(schedule, nil),
		Schedule:    schedule,
	}

	if typed, ok := j.(JobEnabledProvider); ok {
		meta.EnabledProvider = typed.Enabled
	} else {
		meta.EnabledProvider = func() bool { return true }
	}

	jm.jobs[jobName] = meta
	return nil
}

func (jm *JobManager) scheduleNextRuntime(schedule Schedule, after *time.Time) time.Time {
	if schedule != nil {
		return Deref(schedule.GetNextRunTime(after))
	}
	return time.Time{}
}

func (jm *JobManager) setJobDisabledUnsafe(jobName string, disabled bool) error {
	if _, hasJob := jm.jobs[jobName]; !hasJob {
		return exception.New(ErrJobNotLoaded).WithMessagef("job: %s", jobName)
	}

	jm.jobs[jobName].Disabled = disabled
	return nil
}

func (jm *JobManager) createContext() (context.Context, context.CancelFunc) {
	return context.WithCancel(context.Background())
}

// shouldRunJob returns whether it is legal to run a job based off of a job's attributes and status.
// Use this function to set logic for whether a job should run
func (jm *JobManager) shouldRunJob(job Job) bool {
	if meta, ok := jm.jobs[job.Name()]; ok {
		return !meta.Disabled
	}
	return false
}

// shouldRunTask returns whether a task should be executed based on its status
func (jm *JobManager) shouldRunTask(t Task) bool {
	_, isSerial := t.(TaskSerialProvider)
	if isSerial {
		_, hasTask := jm.tasks[t.Name()]
		return !hasTask
	}
	return true
}

func (jm *JobManager) onTaskStart(t Task) {
	if jm.shouldTriggerListeners(t) && jm.log != nil {
		jm.log.Trigger(NewEvent(FlagStarted, t.Name()).WithIsWritable(jm.shouldWriteOutput(t)))
	}

	if receiver, isReceiver := t.(TaskOnStartReceiver); isReceiver {
		receiver.OnStart()
	}
}

func (jm *JobManager) onTaskComplete(t Task, elapsed time.Duration, err error) {
	if jm.shouldTriggerListeners(t) && jm.log != nil {
		flag := FlagComplete
		if err != nil {
			flag = FlagFailed
		}

		jm.log.Trigger(NewEvent(flag, t.Name()).
			WithIsWritable(jm.shouldWriteOutput(t)).
			WithElapsed(elapsed).
			WithErr(err))
	}
	if err != nil {
		jm.err(err)
	}

	if typed, ok := t.(TaskOnCompleteReceiver); ok {
		typed.OnComplete(err)
	}

	// test if this needs to have job lifecycle events fired
	if _, ok := t.(Job); ok {
		if jobMeta, ok := jm.jobs[t.Name()]; ok {
			if err != nil && jobMeta.LastErr == nil {
				if typed, ok := jobMeta.Job.(JobOnBrokenReceiver); ok {
					typed.OnBroken(err)
				}
			}
			if err == nil && jobMeta.LastErr != nil {
				if typed, ok := jobMeta.Job.(JobOnFixedReceiver); ok {
					typed.OnFixed()
				}
			}
			jobMeta.LastElapsed = elapsed
			jobMeta.LastErr = err
		}
	}
}

func (jm *JobManager) onTaskCancelled(t Task, elapsed time.Duration) {
	if jm.shouldTriggerListeners(t) && jm.log != nil {
		jm.log.Trigger(NewEvent(FlagCancelled, t.Name()).
			WithIsWritable(jm.shouldWriteOutput(t)).
			WithElapsed(elapsed))
	}
	if receiver, ok := t.(TaskOnCancellationReceiver); ok {
		receiver.OnCancellation()
	}
}

// ShouldTriggerListeners is a helper function to determine if we should trigger listeners for a given task.
func (jm *JobManager) shouldTriggerListeners(t Task) bool {
	if typed, ok := t.(TaskShouldTriggerListenersProvider); ok {
		return typed.ShouldTriggerListeners()
	}
	return true
}

// ShouldWriteOutput is a helper function to determine if we should write logging output for a task.
func (jm *JobManager) shouldWriteOutput(t Task) bool {
	if typed, ok := t.(TaskShouldWriteOutputProvider); ok {
		return typed.ShouldWriteOutput()
	}
	return true
}

//
// logging helpers
//

func (jm *JobManager) err(err error) {
	if err != nil && jm.log != nil {
		jm.log.Error(err)
	}
}

func (jm *JobManager) fatal(err error) {
	if err != nil && jm.log != nil {
		jm.log.Fatal(err)
	}
}

func (jm *JobManager) errorf(format string, args ...interface{}) {
	if jm.log != nil {
		jm.log.SyncErrorf(format, args...)
	}
}

func (jm *JobManager) debugf(format string, args ...interface{}) {
	if jm.log != nil {
		jm.log.SyncDebugf(format, args...)
	}
}
