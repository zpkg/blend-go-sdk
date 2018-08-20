package cron

import (
	"context"
)

// Job is an interface structs can satisfy to be loaded into the JobManager.
type Job interface {
	Name() string
	Schedule() Schedule
	Execute(ctx context.Context) error
}

// IsJobCancelled check if a job is cancelled
func IsJobCancelled(ctx context.Context) bool {
	select {
	case <-ctx.Done():
		return true
	default:
		return false
	}
}

// EventTriggerListenersProvider is a type that enables or disables logger listeners.
type EventTriggerListenersProvider interface {
	ShouldTriggerListeners() bool
}

// EventShouldWriteOutputProvider is a type that enables or disables logger output for events.
type EventShouldWriteOutputProvider interface {
	ShouldWriteOutput() bool
}

// EnabledProvider is an optional interface that will allow jobs to control if they're enabled.
type EnabledProvider interface {
	Enabled() bool
}
