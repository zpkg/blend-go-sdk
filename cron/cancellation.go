package cron

import "context"

// IsJobCancelled check if a job is cancelled
func IsJobCancelled(ctx context.Context) bool {
	select {
	case <-ctx.Done():
		return true
	default:
		return false
	}
}
