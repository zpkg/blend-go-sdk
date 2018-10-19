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
