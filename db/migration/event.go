package migration

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/blend/go-sdk/ansi"
	"github.com/blend/go-sdk/logger"
)

const (
	// Flag is a logger event flag.
	Flag = "db.migration"

	// FlagStats is a logger event flag.
	FlagStats = "db.migration.stats"
)

// NewEvent returns a new event.
func NewEvent(result, body string, labels ...string) *Event {
	return &Event{
		EventMeta: logger.NewEventMeta(Flag),
		Result:    result,
		Body:      body,
		Labels:    labels,
	}
}

// Event is a migration logger event.
type Event struct {
	*logger.EventMeta

	Result string
	Body   string
	Labels []string
}

func (e Event) colorizeFixedWidthLeftAligned(tf logger.TextFormatter, text string, color ansi.Color, width int) string {
	fixedToken := fmt.Sprintf("%%-%ds", width)
	return tf.Colorize(fmt.Sprintf(fixedToken, text), color)
}

// WriteText writes the migration event as text.
func (e Event) WriteText(tf logger.Colorizer, wr io.Writer) {
	resultColor := ansi.ColorBlue
	switch e.Result {
	case "skipped":
		resultColor = ansi.ColorYellow
	case "failed":
		resultColor = ansi.ColorRed
	}

	if len(e.Result) > 0 {
		io.WriteString(wr, tf.Colorize("--", ansi.ColorLightBlack))
		io.WriteString(wr, logger.Space)
		io.WriteString(wr, tf.Colorize(e.Result, resultColor))
	}

	if len(e.Labels) > 0 {
		io.WriteString(wr, logger.Space)
		io.WriteString(wr, strings.Join(e.Labels, " > "))
	}

	if len(e.Body) > 0 {
		io.WriteString(wr, logger.Space)
		io.WriteString(wr, tf.Colorize("--", ansi.ColorLightBlack))
		io.WriteString(wr, logger.Space)
		io.WriteString(wr, e.Body)
	}
}

// Fields implements logger.FieldsProvider.
func (e Event) Fields() logger.Fields {
	return logger.Fields{
		"result": e.Result,
		"labels": e.Labels,
		"body":   e.Body,
	}
}

// NewStatsEvent returns a new stats event.
func NewStatsEvent(applied, skipped, failed, total int) *StatsEvent {
	return &StatsEvent{
		EventMeta: logger.NewEventMeta(FlagStats),
		applied:   applied,
		skipped:   skipped,
		failed:    failed,
		total:     total,
	}
}

// StatsEvent is a migration logger event.
type StatsEvent struct {
	*logger.EventMeta
	applied int
	skipped int
	failed  int
	total   int
}

// WriteText writes the event to a text writer.
func (se StatsEvent) WriteText(tf logger.TextFormatter, buf *bytes.Buffer) {
	buf.WriteString(fmt.Sprintf("%s applied %s skipped %s failed %s total",
		tf.Colorize(fmt.Sprintf("%d", se.applied), ansi.ColorGreen),
		tf.Colorize(fmt.Sprintf("%d", se.skipped), ansi.ColorLightGreen),
		tf.Colorize(fmt.Sprintf("%d", se.failed), ansi.ColorRed),
		tf.Colorize(fmt.Sprintf("%d", se.total), ansi.ColorLightWhite),
	))
}

// Fields implements logger.FieldsProvider.
func (se StatsEvent) Fields() logger.Fields {
	return logger.Fields{
		StatApplied: se.applied,
		StatSkipped: se.skipped,
		StatFailed:  se.failed,
		StatTotal:   se.total,
	}
}
