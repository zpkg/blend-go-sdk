package jobkit

import (
	"github.com/blend/go-sdk/airbrake"
	"github.com/blend/go-sdk/aws"
	"github.com/blend/go-sdk/cron"
	"github.com/blend/go-sdk/datadog"
	"github.com/blend/go-sdk/email"
	"github.com/blend/go-sdk/logger"
	"github.com/blend/go-sdk/slack"
	"github.com/blend/go-sdk/web"
)

// Config is the jobkit config.
type Config struct {
	cron.Config `json:",inline" yaml:",inline"`

	Logger logger.Config `json:"logger" yaml:"logger"`
	Web    web.Config    `json:"web" yaml:"web"`

	Airbrake airbrake.Config `json:"airbrake" yaml:"airbrake"`
	AWS      aws.Config      `json:"aws" yaml:"aws"`
	Email    email.Message   `json:"email" yaml:"email"`
	Datadog  datadog.Config  `json:"datadog" yaml:"datadog"`
	Slack    slack.Config    `json:"slack" yaml:"slack"`

	Notifications map[string]NotificationsConfig `json:"notifications" yaml:"notifications"`
}
