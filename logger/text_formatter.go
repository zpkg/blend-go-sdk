package logger

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/blend/go-sdk/ansi"
)

// NewTextFormatter returns a new text writer for a given output.
func NewTextFormatter(options ...TextFormatterOption) *TextFormatter {
	tf := &TextFormatter{
		TimeFormat: DefaultTextTimeFormat,
	}

	for _, option := range options {
		option(tf)
	}

	return tf
}

// TextFormatterOption is an option for text formatters.
type TextFormatterOption func(*TextFormatter)

// OptTextConfig sets the text formatter config.
func OptTextConfig(cfg *TextConfig) TextFormatterOption {
	return func(tf *TextFormatter) {
		tf.HideTimestamp = cfg.HideTimestamp
		tf.HideFields = cfg.HideFields
		tf.NoColor = cfg.NoColor
		tf.TimeFormat = cfg.TimeFormatOrDefault()
	}
}

// TextFormatter handles formatting messages as text.
type TextFormatter struct {
	HideTimestamp bool
	HideFields    bool
	NoColor       bool
	TimeFormat    string
}

// Colorize (optionally) applies a color to a string.
func (tf TextFormatter) Colorize(value string, color ansi.Color) string {
	if tf.NoColor {
		return value
	}
	return color.Apply(value)
}

// FormatFlag formats the flag portion of the message.
func (tf TextFormatter) FormatFlag(flag string, color ansi.Color) string {
	return fmt.Sprintf("[%s]", tf.Colorize(string(flag), color))
}

// FormatTimestamp returns a new timestamp string.
func (tf TextFormatter) FormatTimestamp(ts time.Time) string {
	timeFormat := DefaultTextTimeFormat
	if len(tf.TimeFormat) > 0 {
		timeFormat = tf.TimeFormat
	}
	value := ts.Format(timeFormat)
	return tf.Colorize(fmt.Sprintf("%-30s", value), ansi.ColorGray)
}

// FormatSubContextPath returns the sub-context path section of the message.
func (tf *TextFormatter) FormatSubContextPath(path ...string) string {
	if len(path) == 0 {
		return ""
	}
	if len(path) == 1 {
		return fmt.Sprintf("[%s]", tf.Colorize(path[0], ansi.ColorBlue))
	}
	if !tf.NoColor {
		for index := 0; index < len(path); index++ {
			path[index] = tf.Colorize(path[index], ansi.ColorBlue)
		}
	}
	return fmt.Sprintf("[%s]", strings.Join(path, " > "))
}

// WriteFormat implements write formatter.
func (tf TextFormatter) WriteFormat(ctx context.Context, output io.Writer, e Event) error {
	buffer := new(bytes.Buffer)

	if !tf.HideTimestamp {
		buffer.WriteString(tf.FormatTimestamp(e.Timestamp()))
		buffer.WriteString(Space)
	}

	buffer.WriteString(tf.FormatFlag(e.Flag(), FlagTextColor(e.Flag())))
	buffer.WriteString(Space)

	if subContextPath := GetSubContextPath(ctx); subContextPath != nil {
		buffer.WriteString(tf.FormatSubContextPath(subContextPath...))
		buffer.WriteString(Space)
	}

	if typed, ok := e.(TextWritable); ok {
		typed.WriteText(tf, buffer)
	} else if fieldsProvider, ok := e.(FieldsProvider); ok {
		fields := fieldsProvider.Fields()
		WriteFields(tf, buffer, fields)
	} else if stringer, ok := e.(fmt.Stringer); ok {
		buffer.WriteString(stringer.String())
	}

	buffer.WriteString(Newline)
	_, err := io.Copy(output, buffer)
	return err
}
