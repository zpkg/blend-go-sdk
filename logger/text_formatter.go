package logger

import (
	"fmt"
	"io"
	"time"

	"github.com/blend/go-sdk/ansi"
)

// NewTextFormatter returns a new text writer for a given output.
func NewTextFormatter(cfg *TextConfig) *TextFormatter {
	return &TextFormatter{
		HideTimestamp: cfg.HideTimestamp,
		NoColor:       cfg.NoColor,
		TimeFormat:    cfg.TimeFormatOrDefault(),
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
	if wr.NoColor {
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
	if len(wr.timeFormat) > 0 {
		timeFormat = wr.timeFormat
	}
	value := ts.Format(timeFormat)
	return wr.Colorize(fmt.Sprintf("%-30s", value), ansi.ColorGray)
}

// WriteFormat implements write formatter.
func (tf TextFormatter) WriteFormat(output io.Writer, e Event) error {
	if wr.ShowTimestamp {
		buf.WriteString(wr.FormatTimestamp(e.Timestamp()))
		buf.WriteRune(RuneSpace)
	}

	buf.WriteString(wr.FormatFlag(e.Flag(), GetFlagTextColor(e.Flag())))
	buf.WriteRune(RuneSpace)

	buf.WriteString(wr.FormatEntity(typed.Entity(), ansi.ColorBlue))
	buf.WriteRune(RuneSpace)

	e.WriteText(wr, buf)

	buf.WriteRune(RuneNewline)
	_, err := buf.WriteTo(output)
	return err
}
