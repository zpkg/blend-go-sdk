package jobkit

import (
	"context"
	"fmt"
	"time"

	"github.com/blend/go-sdk/cron"
	"github.com/blend/go-sdk/diagnostics"
	"github.com/blend/go-sdk/email"
	"github.com/blend/go-sdk/logger"
	"github.com/blend/go-sdk/slack"
	"github.com/blend/go-sdk/stats"
)

var (
	_ cron.Job                    = (*Job)(nil)
	_ cron.OnStartReceiver        = (*Job)(nil)
	_ cron.OnCompleteReceiver     = (*Job)(nil)
	_ cron.OnFailureReceiver      = (*Job)(nil)
	_ cron.OnCancellationReceiver = (*Job)(nil)
	_ cron.OnBrokenReceiver       = (*Job)(nil)
	_ cron.OnFixedReceiver        = (*Job)(nil)
	_ cron.OnDisabledReceiver     = (*Job)(nil)
	_ cron.OnEnabledReceiver      = (*Job)(nil)
)

// NewJob creates a new exec job.
func NewJob(action func(context.Context) error) *Job {
	return &Job{
		config: &JobConfig{},
		action: action,
	}
}

// Job is the main job body.
type Job struct {
	name   string
	config *JobConfig

	schedule cron.Schedule
	timeout  time.Duration
	action   func(context.Context) error

	log         logger.Log
	statsClient stats.Collector
	slackClient slack.Sender
	emailClient email.Sender
	errorClient diagnostics.Notifier
}

// Name returns the job name.
func (job Job) Name() string {
	if job.name != "" {
		return job.name
	}
	return job.config.NameOrDefault()
}

// WithName sets the name.
func (job *Job) WithName(name string) *Job {
	job.name = name
	return job
}

// Schedule returns the job schedule.
func (job Job) Schedule() cron.Schedule {
	return job.schedule
}

// WithSchedule sets the schedule.
func (job *Job) WithSchedule(schedule cron.Schedule) *Job {
	job.schedule = schedule
	return job
}

// Config returns the job config.
func (job Job) Config() *JobConfig {
	return job.config
}

// WithConfig sets the config.
func (job *Job) WithConfig(cfg *JobConfig) *Job {
	job.config = cfg
	return job
}

// Timeout returns the timeout.
func (job Job) Timeout() time.Duration {
	return job.timeout
}

// WithTimeout sets the job timeout.
func (job *Job) WithTimeout(d time.Duration) *Job {
	job.timeout = d
	return job
}

// WithLogger sets the job logger.
func (job *Job) WithLogger(log logger.Log) *Job {
	job.log = log
	return job
}

// WithStatsClient sets the job stats client.
func (job *Job) WithStatsClient(client stats.Collector) *Job {
	job.statsClient = client
	return job
}

// WithSlackClient sets the job slack client.
func (job *Job) WithSlackClient(client slack.Sender) *Job {
	job.slackClient = client
	return job
}

// WithEmailClient sets the job email client.
func (job *Job) WithEmailClient(client email.Sender) *Job {
	job.emailClient = client
	return job
}

// WithErrorClient sets the job error client.
func (job *Job) WithErrorClient(client diagnostics.Notifier) *Job {
	job.errorClient = client
	return job
}

// OnStart is a lifecycle event handler.
func (job Job) OnStart(ctx context.Context) {
	if job.config != nil && job.config.NotifyOnStartOrDefault() {
		job.notify(ctx, cron.FlagStarted)
	}
}

// OnComplete is a lifecycle event handler.
func (job Job) OnComplete(ctx context.Context) {
	if job.config != nil && job.config.NotifyOnSuccessOrDefault() {
		job.notify(ctx, cron.FlagComplete)
	}
}

// OnFailure is a lifecycle event handler.
func (job Job) OnFailure(ctx context.Context) {
	if job.config != nil && job.config.NotifyOnFailureOrDefault() {
		job.notify(ctx, cron.FlagFailed)
	}
}

// OnBroken is a lifecycle event handler.
func (job Job) OnBroken(ctx context.Context) {
	if job.config != nil && job.config.NotifyOnBrokenOrDefault() {
		job.notify(ctx, cron.FlagBroken)
	}
}

// OnFixed is a lifecycle event handler.
func (job Job) OnFixed(ctx context.Context) {
	if job.config != nil && job.config.NotifyOnFixedOrDefault() {
		job.notify(ctx, cron.FlagFixed)
	}
}

// OnCancellation is a lifecycle event handler.
func (job Job) OnCancellation(ctx context.Context) {
	if job.config != nil && job.config.NotifyOnFailureOrDefault() {
		job.notify(ctx, cron.FlagCancelled)
	}
}

// OnEnabled is a lifecycle event handler.
func (job Job) OnEnabled(ctx context.Context) {
	if job.config != nil && job.config.NotifyOnEnabledOrDefault() {
		job.notify(ctx, cron.FlagEnabled)
	}
}

// OnDisabled is a lifecycle event handler.
func (job Job) OnDisabled(ctx context.Context) {
	if job.config != nil && job.config.NotifyOnDisabledOrDefault() {
		job.notify(ctx, cron.FlagDisabled)
	}
}

func (job Job) notify(ctx context.Context, flag logger.Flag) {
	if job.statsClient != nil {
		job.statsClient.Increment(string(flag), fmt.Sprintf("%s:%s", stats.TagJob, job.Name()))
		if ji := cron.GetJobInvocation(ctx); ji != nil {
			logger.MaybeError(job.log, job.statsClient.TimeInMilliseconds(string(flag), ji.Elapsed, fmt.Sprintf("%s:%s", stats.TagJob, job.Name())))
		}
	}
	if job.slackClient != nil {
		if ji := cron.GetJobInvocation(ctx); ji != nil {
			logger.MaybeError(job.log, job.slackClient.Send(context.Background(), NewSlackMessage(string(flag), ji)))
		}
	}
	if job.emailClient != nil {
		if ji := cron.GetJobInvocation(ctx); ji != nil {
			message, err := NewEmailMessage(string(flag), ji)
			if err != nil {
				logger.MaybeError(job.log, err)
			}
			logger.MaybeError(job.log, job.emailClient.Send(context.Background(), message))
		}
	}
	if job.errorClient != nil {
		if ji := cron.GetJobInvocation(ctx); ji != nil && ji.Err != nil {
			job.errorClient.Notify(ji.Err)
		}
	}
}

// Execute is the job body.
func (job Job) Execute(ctx context.Context) error {
	return job.action(ctx)
}
