package cron

import (
	"context"
	"time"
)

// JobInvocation is metadata for a job invocation (or instance of a job running).
type JobInvocation struct {
	ID        string             `json:"id"`
	Name      string             `json:"name"`
	Started   time.Time          `json:"started"`
	Finished  time.Time          `json:"finished"`
	Cancelled time.Time          `json:"cancelled"`
	Timeout   time.Time          `json:"timeout"`
	Err       error              `json:"err"`
	Elapsed   time.Duration      `json:"elapsed"`
	Status    JobStatus          `json:"status"`
	Context   context.Context    `json:"-"`
	Cancel    context.CancelFunc `json:"-"`
}
