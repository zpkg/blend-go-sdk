package jobkit

import (
	"github.com/blend/go-sdk/logger"
	"github.com/blend/go-sdk/slack"
	"github.com/blend/go-sdk/web"
)

// Config is a jobkit config.
type Config struct {
	Logger logger.Config `json:"logger" yaml:"logger"`
	Web    web.Config    `json:"web" yaml:"web"`
	Slack  slack.Config  `json:"slack" yaml:"slack"`

	Jobs []JobConfig `json:"jobs" yaml:"jobs"`
}
