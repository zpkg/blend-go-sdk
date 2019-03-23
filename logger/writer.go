package logger

import (
	"strings"

	"github.com/blend/go-sdk/env"
)

const (
	// OutputFormatJSON is an output format.
	OutputFormatJSON = "json"
	// OutputFormatText is an output format.
	OutputFormatText = "text"
)

// NewWriter creates a new writer based on a given format.
// It reads the writer settings from the environment.
func NewWriter(format string) Writer {
	switch strings.ToLower(string(format)) {
	case OutputFormatJSON:
		return NewJSONWriterFromEnv()
	case OutputFormatText:
		return NewTextWriterFromEnv()
	}
	panic("invalid writer output format")
}

// NewWriterFromEnv returns a new writer based on the environment variable `LOG_FORMAT`.
func NewWriterFromEnv() Writer {
	if format := env.Env().String(EnvVarFormat); len(format) > 0 {
		return NewWriter(OutputFormat(format))
	}
	return NewTextWriterFromEnv()
}
