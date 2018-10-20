package cron

import "time"

/*
A note on the naming conventions for these interfaces.

MethodName[Receiver|Provider] is the general pattern.

If possible, Recever indicates a method that will be called actively, provider is a value your type can give.

They're mostly the same except
*/

// TimeoutProvider is an interface that allows a task to be timed out.
type TimeoutProvider interface {
	Timeout() time.Duration
}

// StatusProvider is an interface that allows a task to report its status.
type StatusProvider interface {
	Status() string
}

// SerialProvider is an optional interface that prohibits
// a task from running if another instance of the task is currently running.
type SerialProvider interface {
	Serial()
}

// ShouldTriggerListenersProvider is a type that enables or disables logger listeners.
type ShouldTriggerListenersProvider interface {
	ShouldTriggerListeners() bool
}

// ShouldWriteOutputProvider is a type that enables or disables logger output for events.
type ShouldWriteOutputProvider interface {
	ShouldWriteOutput() bool
}

// JobEnabledProvider is an optional interface that will allow jobs to control if they're enabled.
type JobEnabledProvider interface {
	Enabled() bool
}

// OnStartReceiver is an interface that allows a task to be signaled when it has started.
type OnStartReceiver interface {
	OnStart(*TaskInvocation)
}

// OnCancellationReceiver is an interface that allows a task to be signaled when it has been canceled.
type OnCancellationReceiver interface {
	OnCancellation(*TaskInvocation)
}

// OnCompleteReceiver is an interface that allows a task to be signaled when it has been completed.
type OnCompleteReceiver interface {
	OnComplete(*TaskInvocation)
}

// JobOnBrokenReceiver is an interface that allows a job to be signaled when it is a failure that followed
// a previous success.
type JobOnBrokenReceiver interface {
	OnBroken(*TaskInvocation)
}

// JobOnFixedReceiver is an interface that allows a jbo to be signaled when is a success that followed
// a previous failure.
type JobOnFixedReceiver interface {
	OnFixed(*TaskInvocation)
}
