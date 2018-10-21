package cron

import (
	"context"
	"time"
)

// JobInvocation is metadata for a job invocation (or instance of a job running).
type JobInvocation struct {
	Name      string             `json:"name"`
	JobMeta   *JobMeta           `json:"jobMeta"`
	StartTime time.Time          `json:"startTime"`
	Timeout   time.Time          `json:"timeout"`
	Context   context.Context    `json:"-"`
	Cancel    context.CancelFunc `json:"-"`
	Err       error              `json:"err"`
	Elapsed   time.Duration      `json:"elapsed"`
}
