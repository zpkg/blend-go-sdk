package jobkit

import (
	"context"

	"github.com/blend/go-sdk/airbrake"
	"github.com/blend/go-sdk/aws"
	"github.com/blend/go-sdk/aws/ses"
	"github.com/blend/go-sdk/cron"
	"github.com/blend/go-sdk/datadog"
	"github.com/blend/go-sdk/diagnostics"
	"github.com/blend/go-sdk/email"
	"github.com/blend/go-sdk/slack"
	"github.com/blend/go-sdk/stats"
)

// New returns a new job.
func New(jobConfig *JobConfig, cfg *Config, action func(context.Context) error) (*Job, error) {
	schedule, err := cron.ParseString(jobConfig.ScheduleOrDefault())
	if err != nil {
		return nil, err
	}

	// set up myriad of notification targets
	var emailClient email.Sender
	if !cfg.AWS.IsZero() {
		emailClient = ses.New(aws.MustNewSession(&cfg.AWS))
	}
	var slackClient slack.Sender
	if !cfg.Slack.IsZero() {
		slackClient = slack.New(&cfg.Slack)
	}
	var statsClient stats.Collector
	if !cfg.Datadog.IsZero() {
		statsClient, err = datadog.NewCollector(&cfg.Datadog)
		if err != nil {
			return nil, err
		}
	}
	var errorClient diagnostics.Notifier
	if !cfg.Airbrake.IsZero() {
		errorClient = airbrake.MustNew(&cfg.Airbrake)
	}

	job := NewJob(action).
		WithName(jobConfig.NameOrDefault()).
		WithConfig(jobConfig).
		WithSchedule(schedule).
		WithEmailClient(emailClient).
		WithStatsClient(statsClient).
		WithSlackClient(slackClient).
		WithErrorClient(errorClient)

	return job, nil
}
