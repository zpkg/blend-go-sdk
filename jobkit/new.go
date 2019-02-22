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

// NewJob returns a new job.
func NewJob(jobConfig *JobConfig, kitConfig *Config, action func(context.Context) error) (*Job, error) {
	schedule, err := cron.ParseString(jobConfig.ScheduleOrDefault())
	if err != nil {
		return nil, err
	}

	// set up myriad of notification targets
	var emailClient email.Sender
	if !kitConfig.AWS.IsZero() {
		emailClient = ses.New(aws.MustNewSession(&kitConfig.AWS))
	}
	var slackClient slack.Sender
	if !kitConfig.Slack.IsZero() {
		slackClient = slack.New(&kitConfig.Slack)
	}
	var statsClient stats.Collector
	if !kitConfig.Datadog.IsZero() {
		statsClient, err = datadog.NewCollector(&kitConfig.Datadog)
		if err != nil {
			return nil, err
		}
	}
	var errorClient diagnostics.Notifier
	if !kitConfig.Airbrake.IsZero() {
		errorClient = airbrake.MustNew(&kitConfig.Airbrake)
	}

	job := (&Job{action: action}).
		WithName(jobConfig.NameOrDefault()).
		WithConfig(jobConfig).
		WithSchedule(schedule).
		WithTimeout(jobConfig.TimeoutOrDefault()).
		WithEmailClient(emailClient).
		WithStatsClient(statsClient).
		WithSlackClient(slackClient).
		WithErrorClient(errorClient)

	return job, nil
}
