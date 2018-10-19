package cron

import "time"

// TaskResumeProvider is an interface that allows a task to be resumed.
type TaskResumeProvider interface {
	State() interface{}
	Resume(state interface{}) error
}

// TaskTimeoutProvider is an interface that allows a task to be timed out.
type TaskTimeoutProvider interface {
	Timeout() time.Duration
}

// TaskStatusProvider is an interface that allows a task to report its status.
type TaskStatusProvider interface {
	Status() string
}

// TaskOnStartReceiver is an interface that allows a task to be signaled when it has started.
type TaskOnStartReceiver interface {
	OnStart()
}

// TaskOnCancellationReceiver is an interface that allows a task to be signaled when it has been canceled.
type TaskOnCancellationReceiver interface {
	OnCancellation()
}

// TaskOnCompleteReceiver is an interface that allows a task to be signaled when it has been completed.
type TaskOnCompleteReceiver interface {
	OnComplete(error)
}

// TaskSerialProvider is an optional interface that prohibits
// a task from running if another instance of the task is currently running.
type TaskSerialProvider interface {
	Serial()
}

// TaskShouldTriggerListenersProvider is a type that enables or disables logger listeners.
type TaskShouldTriggerListenersProvider interface {
	ShouldTriggerListeners() bool
}

// TaskShouldWriteOutputProvider is a type that enables or disables logger output for events.
type TaskShouldWriteOutputProvider interface {
	ShouldWriteOutput() bool
}

// JobEnabledProvider is an optional interface that will allow jobs to control if they're enabled.
type JobEnabledProvider interface {
	Enabled() bool
}

// JobOnBrokenReceiver is an interface that allows a job to be signaled when it is a failure that followed
// a previous success.
type JobOnBrokenReceiver interface {
	OnBroken(error)
}

// JobOnFixedReceiver is an interface that allows a jbo to be signaled when is a success that followed
// a previous failure.
type JobOnFixedReceiver interface {
	OnFixed()
}
