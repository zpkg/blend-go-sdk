package cron

import (
	"context"
	"time"
)

// TaskInvocation is metadata for a running task.
type TaskInvocation struct {
	Name      string             `json:"name"`
	Task      Task               `json:"-"`
	StartTime time.Time          `json:"startTime"`
	Timeout   time.Time          `json:"timeout"`
	Context   context.Context    `json:"-"`
	Cancel    context.CancelFunc `json:"-"`
	Err       error              `json:"err"`
	Elapsed   time.Duration      `json:"elapsed"`
}
