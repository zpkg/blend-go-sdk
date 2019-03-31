package jobkit

import (
	"context"

	"github.com/blend/go-sdk/cron"
	"github.com/blend/go-sdk/logger"
)

// LogSubContextDebugf is a logger interface.
type LogSubContextDebugf interface {
	logger.DebugfReceiver
	logger.SubContexter
}

// LogSubContextInfof is a logger interface.
type LogSubContextInfof interface {
	logger.InfofReceiver
	logger.SubContexter
}

// LogSubContextWarningf is a logger interface.
type LogSubContextWarningf interface {
	logger.WarningfReceiver
	logger.SubContexter
}

// LogSubContextWarning is a logger interface.
type LogSubContextWarning interface {
	logger.WarningReceiver
	logger.SubContexter
}

// LogSubContextErrorf is a logger interface.
type LogSubContextErrorf interface {
	logger.ErrorfReceiver
	logger.SubContexter
}

// LogSubContextError is a logger interface.
type LogSubContextError interface {
	logger.ErrorReceiver
	logger.SubContexter
}

// LogSubContextFatalf is a logger interface.
type LogSubContextFatalf interface {
	logger.FatalfReceiver
	logger.SubContexter
}

// LogSubContextFatal is a logger interface.
type LogSubContextFatal interface {
	logger.FatalReceiver
	logger.SubContexter
}

// Debugf prints an info message if the logger is set.
func Debugf(ctx context.Context, log LogSubContextDebugf, format string, args ...interface{}) {
	if log == nil {
		return
	}
	ji := cron.GetJobInvocation(ctx)
	log.SubContext(ji.ID).Debugf(format, args...)
}

// Infof prints an info message if the logger is set.
func Infof(ctx context.Context, log LogSubContextInfof, format string, args ...interface{}) {
	if log == nil {
		return
	}
	ji := cron.GetJobInvocation(ctx)
	log.SubContext(ji.ID).Infof(format, args...)
}

// Warningf prints a warning message if the logger is set.
func Warningf(ctx context.Context, log LogSubContextWarningf, format string, args ...interface{}) {
	if log == nil {
		return
	}
	ji := cron.GetJobInvocation(ctx)
	log.SubContext(ji.ID).Warningf(format, args...)
}

// Warning prints an warning if the logger is set.
func Warning(ctx context.Context, log LogSubContextWarning, err error) {
	if log == nil {
		return
	}
	ji := cron.GetJobInvocation(ctx)
	log.SubContext(ji.ID).Warning(err)
}

// Errorf prints an error message if the logger is set.
func Errorf(ctx context.Context, log LogSubContextErrorf, format string, args ...interface{}) {
	if log == nil {
		return
	}
	ji := cron.GetJobInvocation(ctx)
	log.SubContext(ji.ID).Errorf(format, args...)
}

// Error prints an error if the logger is set.
func Error(ctx context.Context, log LogSubContextError, err error) {
	if log == nil {
		return
	}
	ji := cron.GetJobInvocation(ctx)
	log.SubContext(ji.ID).Error(err)
}

// Fatalf prints a fatal error message if the logger is set.
func Fatalf(ctx context.Context, log LogSubContextFatalf, format string, args ...interface{}) {
	if log == nil {
		return
	}
	ji := cron.GetJobInvocation(ctx)
	log.SubContext(ji.ID).Fatalf(format, args...)
}

// Fatal prints a fatal error if the logger is set.
func Fatal(ctx context.Context, log LogSubContextFatal, err error) {
	if log == nil {
		return
	}
	ji := cron.GetJobInvocation(ctx)
	log.SubContext(ji.ID).Fatal(err)
}
