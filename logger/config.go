package logger

import (
	"strings"
)

// Config is the logger config.
type Config struct {
	Flags  []string   `json:"flags,omitempty" yaml:"flags,omitempty" env:"LOG_FLAGS,csv"`
	Format string     `json:"format,omitempty" yaml:"format,omitempty" env:"LOG_FORMAT"`
	Text   TextConfig `json:"text,omitempty" yaml:"text,omitempty"`
	JSON   JSONConfig `json:"json,omitempty" yaml:"json,omitempty"`
}

// FlagsOrDefault returns the enabled logger events.
func (c Config) FlagsOrDefault() []string {
	if len(c.Flags) > 0 {
		return c.Flags
	}
	return DefaultFlags
}

// FormatOrDefault returns the output format or a default.
func (c Config) FormatOrDefault() string {
	if c.Format != "" {
		return c.Format
	}
	return FormatText
}

// Formatter returns the configured writers
func (c Config) Formatter() WriteFormatter {
	switch strings.ToLower(string(c.FormatOrDefault())) {
	case FormatJSON:
		return NewJSONFormatter(&c.JSON)
	case FormatText:
		return NewTextFormatter(&c.Text)
	default:
		return NewTextFormatter(&c.Text)
	}
}

// TextConfig is the config for a text formatter.
type TextConfig struct {
	HideTimestamp bool   `json:"hideTimestamp,omitempty" yaml:"hideTimestamp,omitempty" env:"LOG_HIDE_TIMESTAMP"`
	NoColor       bool   `json:"noColor,omitempty" yaml:"noColor,omitempty" env:"NO_COLOR"`
	TimeFormat    string `json:"timeFormat,omitempty" yaml:"timeFormat,omitempty" env:"LOG_TIME_FORMAT"`
}

// TimeFormatOrDefault returns a field value or a default.
func (twc TextConfig) TimeFormatOrDefault() string {
	if len(twc.TimeFormat) > 0 {
		return twc.TimeFormat
	}
	return DefaultTextTimeFormat
}

// JSONConfig is the config for a json formatter.
type JSONConfig struct {
	Pretty bool `json:"pretty,omitempty" yaml:"pretty,omitempty" env:"LOG_JSON_PRETTY"`
}
