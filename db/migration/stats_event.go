package migration

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/blend/go-sdk/ansi"
	"github.com/blend/go-sdk/logger"
)

var (
	_ logger.Event        = (*StatsEvent)(nil)
	_ logger.TextWritable = (*StatsEvent)(nil)
	_ json.Marshaler      = (*StatsEvent)(nil)
)

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
func (se StatsEvent) WriteText(tf logger.TextFormatter, wr io.Writer) {
	io.WriteString(wr, fmt.Sprintf("%s applied %s skipped %s failed %s total",
		tf.Colorize(fmt.Sprintf("%d", se.applied), ansi.ColorGreen),
		tf.Colorize(fmt.Sprintf("%d", se.skipped), ansi.ColorLightGreen),
		tf.Colorize(fmt.Sprintf("%d", se.failed), ansi.ColorRed),
		tf.Colorize(fmt.Sprintf("%d", se.total), ansi.ColorLightWhite),
	))
}

// MarshalJSON implements json.Marshaler.
func (se StatsEvent) MarshalJSON() ([]byte, error) {
	return json.Marshal(logger.MergeDecomposed(se.EventMeta.Decompose(), map[string]interface{}{
		StatApplied: se.applied,
		StatSkipped: se.skipped,
		StatFailed:  se.failed,
		StatTotal:   se.total,
	}))
}
