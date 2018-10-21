package cron

import (
	"context"
	"time"
)

var (
	_ ScheduleProvider               = (*JobBuilder)(nil)
	_ TimeoutProvider                = (*JobBuilder)(nil)
	_ EnabledProvider                = (*JobBuilder)(nil)
	_ ShouldWriteOutputProvider      = (*JobBuilder)(nil)
	_ ShouldTriggerListenersProvider = (*JobBuilder)(nil)
	_ OnStartReceiver                = (*JobBuilder)(nil)
	_ OnCancellationReceiver         = (*JobBuilder)(nil)
	_ OnCompleteReceiver             = (*JobBuilder)(nil)
	_ OnFailureReceiver              = (*JobBuilder)(nil)
	_ OnBrokenReceiver               = (*JobBuilder)(nil)
	_ OnFixedReceiver                = (*JobBuilder)(nil)
)

// NewJob returns a new job factory.
func NewJob(name string) *JobBuilder {
	return &JobBuilder{
		name: name,
	}
}

// JobBuilder allows for job creation w/o a fully formed struct.
type JobBuilder struct {
	name                 string
	timeoutProvider      func() time.Duration
	enabledProvider      func() bool
	showMessagesProvider func() bool
	schedule             Schedule
	action               TaskAction

	onStart     func(*TaskInvocation)
	onCancelled func(*TaskInvocation)
	onComplete  func(*TaskInvocation)
	onFailure   func(*TaskInvocation)
	onBroken    func(*TaskInvocation)
	onFixed     func(*TaskInvocation)
}

// WithName sets the job name.
func (jb *JobBuilder) WithName(name string) *JobBuilder {
	jb.name = name
	return jb
}

// WithSchedule sets the schedule for the job.
func (jb *JobBuilder) WithSchedule(schedule Schedule) *JobBuilder {
	jb.schedule = schedule
	return jb
}

// TimeoutProvider returns the job timeout.
func (jb *JobBuilder) TimeoutProvider() func() time.Duration {
	return jb.timeoutProvider
}

// WithTimeoutProvider sets the timeout provider.
func (jb *JobBuilder) WithTimeoutProvider(timeoutProvider func() time.Duration) *JobBuilder {
	jb.timeoutProvider = timeoutProvider
	return jb
}

// WithAction sets the job action.
func (jb *JobBuilder) WithAction(action TaskAction) *JobBuilder {
	jb.action = action
	return jb
}

// Action returns the job action.
func (jb *JobBuilder) Action() TaskAction {
	return jb.action
}

// WithEnabledProvider sets the enabled provider for the job.
func (jb *JobBuilder) WithEnabledProvider(enabledProvider func() bool) *JobBuilder {
	jb.enabledProvider = enabledProvider
	return jb
}

// WithShowMessagesProvider sets the enabled provider for the job.
func (jb *JobBuilder) WithShowMessagesProvider(provider func() bool) *JobBuilder {
	jb.showMessagesProvider = provider
	return jb
}

//
// implementations of interface methods
//

// Name returns the job name.
func (jb *JobBuilder) Name() string {
	return jb.name
}

// Schedule returns the job schedule.
func (jb *JobBuilder) Schedule() Schedule {
	return jb.schedule
}

// Enabled returns if the job is enabled.
func (jb *JobBuilder) Enabled() bool {
	if jb.enabledProvider != nil {
		return jb.enabledProvider()
	}
	return true
}

// Execute runs the job action if it's set.
func (jb *JobBuilder) Execute(ctx context.Context) error {
	if jb.action != nil {
		return jb.action(ctx)
	}
	return nil
}

// ShowMessages returns if the job should trigger logging events.
func (jb *JobBuilder) ShowMessages() bool {
	if jb.showMessagesProvider != nil {
		return jb.showMessagesProvider()
	}
	return true
}
