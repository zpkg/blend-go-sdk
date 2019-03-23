package logger

import (
	"strings"

	"github.com/blend/go-sdk/env"
)

// Config is the logger config.
type Config struct {
	Flags        []string         `json:"flags,omitempty" yaml:"flags,omitempty" env:"LOG_EVENTS,csv"`
	OutputFormat string           `json:"outputFormat,omitempty" yaml:"outputFormat,omitempty" env:"LOG_OUTPUT_FORMAT"`
	Text         TextWriterConfig `json:"text,omitempty" yaml:"text,omitempty"`
	JSON         JSONWriterConfig `json:"json,omitempty" yaml:"json,omitempty"`
}

// FlagsOrDefault returns the enabled logger events.
func (c Config) FlagsOrDefault() []string {
	if len(c.Flags) > 0 {
		return c.Flags
	}
	return DefaultFlags
}

// OutputFormatOrDefault returns the output format or a default.
func (c Config) OutputFormatOrDefault() string {
	if c.OutputFormat != "" {
		return c.OutputFormat
	}
	return OutputFormatText
}

// Writers returns the configured writers
func (c Config) Writers() []Writer {
	switch OutputFormat(strings.ToLower(string(c.OutputFormatOrDefault()))) {
	case OutputFormatJSON:
		return []Writer{NewJSONWriterFromConfig(&c.JSON)}
	case OutputFormatText:
		return []Writer{NewTextWriterFromConfig(&c.Text)}
	default:
		return []Writer{NewTextWriterFromConfig(&c.Text)}
	}
}

// NewTextWriterConfigFromEnv returns a new text writer config from the environment.
func NewTextWriterConfigFromEnv() *TextWriterConfig {
	var config TextWriterConfig
	if err := env.Env().ReadInto(&config); err != nil {
		panic(err)
	}
	return &config
}

// TextWriterConfig is the config for a text writer.
type TextWriterConfig struct {
	HideTimestamp bool   `json:"hideTimestamp,omitempty" yaml:"hideTimestamp,omitempty" env:"LOG_HIDE_TIMESTAMP"`
	NoColor       bool   `json:"noColor,omitempty" yaml:"noColor,omitempty" env:"NO_COLOR"`
	TimeFormat    string `json:"timeFormat,omitempty" yaml:"timeFormat,omitempty" env:"LOG_TIME_FORMAT"`
}

// TimeFormatOrDefault returns a field value or a default.
func (twc TextWriterConfig) TimeFormatOrDefault(defaults ...string) string {
	if len(twc.TimeFormat) > 0 {
		return twc.TimeFormat
	}
	if len(defaults) > 0 {
		return defaults[0]
	}
	return DefaultTextTimeFormat
}

// JSONWriterConfig is the config for a json writer.
type JSONWriterConfig struct {
	Pretty bool `json:"pretty,omitempty" yaml:"pretty,omitempty" env:"LOG_JSON_PRETTY"`
}
