package logger

import (
	"io"
	"os"

	"github.com/blend/go-sdk/env"
)

// Option is a logger option.
type Option func(*Logger)

// OptConfig sets the logger based on a config.
func OptConfig(cfg *Config) Option {
	return func(l *Logger) {
		l.Output = NewInterlockedWriter(os.Stdout)
		l.Formatter = cfg.Formatter()
		l.Flags = NewFlags(cfg.FlagsOrDefault()...)
	}
}

// OptMustConfigFromEnv sets the logger based on a config read from the environment.
// It will panic if there is an erro.
func OptMustConfigFromEnv() Option {
	return func(l *Logger) {
		cfg := &Config{}
		if err := env.Env().ReadInto(cfg); err != nil {
			panic(err)
		}
		l.Output = NewInterlockedWriter(os.Stdout)
		l.Formatter = cfg.Formatter()
		l.Flags = NewFlags(cfg.FlagsOrDefault()...)
	}
}

// OptOutput sets the output writer for the logger.
// It will wrap the output with a synchronizer if it's not already wrapped.
// You can also use this option to "unset" the output by passing in nil.
func OptOutput(output io.Writer) Option {
	return func(l *Logger) {
		if output != nil {
			l.Output = NewInterlockedWriter(output)
		} else {
			l.Output = nil
		}
	}
}

// OptJSON sets the output formatter for the logger as json.
func OptJSON() Option {
	return func(l *Logger) { l.Formatter = NewJSONOutputFormatter() }
}

// OptText sets the output formatter for the logger as json.
func OptText() Option {
	return func(l *Logger) { l.Formatter = NewTextOutputFormatter() }
}

// OptFormatter sets the output formatter.
func OptFormatter(formatter WriteFormatter) Option {
	return func(l *Logger) { l.Formatter = formatter }
}

// OptFlags sets the flags on the logger.
func OptFlags(flags *Flags) Option {
	return func(l *Logger) { l.Flags = flags }
}

// OptAll sets all flags enabled on the logger by default.
func OptAll() Option {
	return func(l *Logger) { l.Flags.SetAll() }
}

// OptNone sets no flags enabled on the logger by default.
func OptNone() Option {
	return func(l *Logger) { l.Flags.SetNone() }
}

// OptEnabled sets enabled flags on the logger.
func OptEnabled(flags ...string) Option {
	return func(l *Logger) { l.Flags.Enable(flags...) }
}

// OptDisabled sets disabled flags on the logger.
func OptDisabled(flags ...string) Option {
	return func(l *Logger) { l.Flags.Disable(flags...) }
}
