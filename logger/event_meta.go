package logger

import (
	"time"

	"github.com/blend/go-sdk/ansi"
)

// these are compile time assertions
var (
	_ Event = (*EventMeta)(nil)
)

// NewEventMeta returns a new event meta.
func NewEventMeta(flag string, options ...EventMetaOption) *EventMeta {
	em := &EventMeta{
		flag:      flag,
		timestamp: time.Now().UTC(),
	}
	for _, option := range options {
		option(em)
	}
	return em
}

// EventMetaOption is an option for event metas.
type EventMetaOption func(*EventMeta)

// OptEventMetaFlag sets the event flag.
func OptEventMetaFlag(flag string) EventMetaOption {
	return func(em *EventMeta) { em.flag = flag }
}

// OptEventMetaTimestamp sets the event timestamp.
func OptEventMetaTimestamp(ts time.Time) EventMetaOption {
	return func(em *EventMeta) { em.timestamp = ts }
}

// OptEventMetaFlagColor sets the event flag color.
func OptEventMetaFlagColor(color ansi.Color) EventMetaOption {
	return func(em *EventMeta) { em.flagColor = color }
}

// EventMeta is the metadata common to events.
// It is useful for ensuring you have the minimum required fields on your events, and its typically embedded in types.
type EventMeta struct {
	flag      string
	timestamp time.Time
	flagColor ansi.Color
}

// Flag returns the event flag.
func (em EventMeta) Flag() string { return em.flag }

// Timestamp returns the event timestamp.
func (em EventMeta) Timestamp() time.Time { return em.timestamp }

// FlagColor returns the event flag color
func (em EventMeta) FlagColor() ansi.Color { return em.flagColor }

// Decompose decomposes the object into a map[string]interface{}.
func (em EventMeta) Decompose() map[string]interface{} {
	output := map[string]interface{}{
		FieldFlag:      em.flag,
		FieldTimestamp: em.timestamp.Format(time.RFC3339Nano),
	}
	return output
}
