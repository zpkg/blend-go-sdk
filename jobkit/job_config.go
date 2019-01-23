package jobkit

import (
	"time"

	"github.com/blend/go-sdk/aws/ses"
	"github.com/blend/go-sdk/slack"
	"github.com/blend/go-sdk/webutil"
)

// JobConfig is something you can use to give your jobs some knobs to turn
// from configuration.
// You can use this job config by embedding it into your larger job config struct.
type JobConfig struct {
	// Name is the name of the job.
	Name string `json:"name" yaml:"name"`
	// Schedule should be the job schedule as represented as a cron string.
	// It should supercede your jobs built in default schedule.
	Schedule string `json:"schedule" yaml:"schedule"`
	// Timeout is an optional timeout to reference for your job.
	Timeout time.Duration `json:"timeout" yaml:"timeout"`

	// NotifyOnBroken governs if we should send notifications on a success => failure transition.
	NotifyOnBroken *bool `json:"notifyOnBroken" yaml:"notifyOnBroken"`
	// NotifyOnFixed governs if we should send notifications on a failure => success transition.
	NotifyOnFixed *bool `json:"notifyOnFixed" yaml:"notifyOnFixed"`
	// NotifyOnSuccess governs if we should send notifications on any success.
	NotifyOnSuccess *bool `json:"notifyOnSuccess" yaml:"notifyOnSuccess"`
	// NotifyOnFailure governs if we should send notifications on any failure.
	NotifyOnFailure *bool `json:"notifyOnFailure" yaml:"notifyOnFailure"`

	// Slack governs slack notification options.
	Slack slack.Config `json:"slack" yaml:"slack"`
	// Webhook governs webhook notification options.
	Webhook webutil.Webhook `json:"webhook" yaml:"webhook"`
	// Email controls email notification defaults.
	Email ses.Message `json:"email" yaml:"email"`
}
