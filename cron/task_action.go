package cron

import "context"

// TaskAction is an function that can be run as a task
type TaskAction func(ctx context.Context) error
