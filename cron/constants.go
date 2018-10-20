package cron

import (
	"time"

	"github.com/blend/go-sdk/logger"
)

const (
	// EnvVarHeartbeatInterval is an environment variable name.
	EnvVarHeartbeatInterval = "CRON_HEARTBEAT_INTERVAL"
)

const (
	// DefaultHeartbeatInterval is the interval between schedule next run checks.
	DefaultHeartbeatInterval = 100 * time.Millisecond

	// DefaultHighPrecisionHeartbeatInterval is the high precision interval between schedule next run checks.
	DefaultHighPrecisionHeartbeatInterval = 10 * time.Millisecond
)

const (
	// FlagStarted is an event flag.
	FlagStarted logger.Flag = "cron.started"
	// FlagFailed is an event flag.
	FlagFailed logger.Flag = "cron.failed"
	// FlagCancelled is an event flag.
	FlagCancelled logger.Flag = "cron.cancelled"
	// FlagComplete is an event flag.
	FlagComplete logger.Flag = "cron.complete"
)

// State is a job state.
type State string

const (
	//StateRunning is the running state.
	StateRunning State = "running"

	// StateEnabled is the enabled state.
	StateEnabled State = "enabled"

	// StateDisabled is the disabled state.
	StateDisabled State = "disabled"
)
