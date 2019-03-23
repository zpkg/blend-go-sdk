package logger

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/blend/go-sdk/ansi"
)

// Asserts text writer is a writer.
var (
	_ Writer = &TextWriter{}
)

// NewTextWriter returns a new text writer for a given output.
func NewTextWriter(output io.Writer) *TextWriter {
	return &TextWriter{
		Output:        NewInterlockedWriter(output),
		ShowHeadings:  DefaultTextWriterShowHeadings,
		ShowTimestamp: DefaultTextWriterShowTimestamp,
		useColor:      DefaultTextWriterUseColor,
		TimeFormat:    DefaultTextTimeFormat,
	}
}

// TextWriter handles outputting logging events to given writer streams as textual columns.
type TextWriter struct {
	Output        io.Writer
	ErrorOutput   io.Writer
	ShowTimestamp bool
	ShowHeadings  bool
	UseColor      bool
	TimeFormat    string
}

// Colorize (optionally) applies a color to a string.
func (wr *TextWriter) Colorize(value string, color ansi.Color) string {
	if wr.UseColor {
		return color.Apply(value)
	}
	return value
}

// ColorizeStatusCode adds color to a status code.
func (wr *TextWriter) ColorizeStatusCode(statusCode int) string {
	if wr.UseColor {
		return ColorizeStatusCode(statusCode)
	}
	return strconv.Itoa(statusCode)
}

// ColorizeByStatusCode colorizes a string by a status code (green, yellow, red).
func (wr *TextWriter) ColorizeByStatusCode(statusCode int, value string) string {
	if wr.UseColor {
		return ColorizeByStatusCode(statusCode, value)
	}
	return value
}

// FormatFlag formats the flag portion of the message.
func (wr *TextWriter) FormatFlag(flag Flag, color ansi.Color) string {
	return fmt.Sprintf("[%s]", wr.Colorize(string(flag), color))
}

// FormatEntity formats the flag portion of the message.
func (wr *TextWriter) FormatEntity(entity string, color ansi.Color) string {
	return fmt.Sprintf("[%s]", wr.Colorize(entity, color))
}

// FormatHeadings returns the headings section of the message.
func (wr *TextWriter) FormatHeadings(headings ...string) string {
	if len(headings) == 0 {
		return ""
	}
	if len(headings) == 1 {
		return fmt.Sprintf("[%s]", wr.Colorize(headings[0], ansi.ColorBlue))
	}
	if wr.useColor {
		for index := 0; index < len(headings); index++ {
			headings[index] = wr.Colorize(headings[index], ansi.ColorBlue)
		}
	}
	return fmt.Sprintf("[%s]", strings.Join(headings, " > "))
}

// FormatTimestamp returns a new timestamp string.
func (wr *TextWriter) FormatTimestamp(ts time.Time) string {
	timeFormat := DefaultTextTimeFormat
	if len(wr.timeFormat) > 0 {
		timeFormat = wr.timeFormat
	}
	value := ts.Format(timeFormat)
	return wr.Colorize(fmt.Sprintf("%-30s", value), ansi.ColorGray)
}

// Write writes to stdout.
func (wr *TextWriter) Write(e Event) error {
	return wr.WriteOutput(wr.Output, e)
}

// WriteError writes to stderr (or stdout if .errorOutput is unset).
func (wr *TextWriter) WriteError(e Event) error {
	if wr.ErrorOutput != nil {
		return wr.WriteOutput(wr.ErrorOutput, e)
	}
	return w.WriteOutput(wr.Output, e)
}

// WriteOutput writes an event to a given output.
func (wr *TextWriter) WriteOutput(output io.Writer, e Event) error {
	buf := new(bytes.Buffer)

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
