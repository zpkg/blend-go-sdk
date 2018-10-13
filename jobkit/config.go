package jobkit

import (
	"github.com/blend/go-sdk/slack"
	"github.com/blend/go-sdk/stats/datadog"
)

// Config are all the options you might need to set for a job.
type Config struct {
	JobName string `json:"jobName" yaml:"jobName" env:"JOB_NAME"`
	JobEnv  string `json:"jobEnv" yaml:"jobEnv" env:"JOB_ENV"`

	Datadog datadog.Config `json:"datadog" yaml:"datadog"`
	Slack   slack.Config   `json:"slack" yaml:"slack"`
}
