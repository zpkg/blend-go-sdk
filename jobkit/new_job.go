package jobkit

import (
	"context"

	"github.com/blend/go-sdk/cron"
)

// NewJob returns a new job.
func NewJob(cfg JobConfig, action func(context.Context) error) (*Job, error) {
	schedule, err := cron.ParseString(cfg.ScheduleOrDefault())
	if err != nil {
		return nil, err
	}

	job := (&Job{action: action}).
		WithName(cfg.Name).
		WithDescription(cfg.Description).
		WithConfig(cfg).
		WithSchedule(schedule).
		WithTimeout(cfg.Timeout)

	return job, nil
}
