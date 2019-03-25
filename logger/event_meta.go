package logger

import (
	"time"

	"github.com/blend/go-sdk/ansi"
)

// these are compile time assertions
var (
	_ Event = &EventMeta{}
)

// NewEventMeta returns a new event meta.
func NewEventMeta(flag string, options ...EventMetaOption) *EventMeta {
	em := &EventMeta{
		flag: flag,
		ts:   time.Now().UTC(),
	}
	for _, option := range options {
		option(em)
	}
	return em
}

// EventMetaOption is an option for event metas.
type EventMetaOption func(*EventMeta)

// EventMetaFlagColor sets the event flag color.
func EventMetaFlagColor(color ansi.Color) EventMetaOption {
	return func(em *EventMeta) { em.flagColor = color }
}

// EventMeta is the metadata common to events.
// It is useful for ensuring you have the minimum required fields on your events, and its typically embedded in types.
type EventMeta struct {
	flag      string
	ts        time.Time
	flagColor ansi.Color

	Labels      map[string]string
	Annotations map[string]string
}

// Flag returns the event flag.
func (em EventMeta) Flag() string { return em.flag }

// FlagColor returns the event flag color
func (em EventMeta) FlagColor() ansi.Color { return em.flagColor }

// Timestamp returns the event timestamp.
func (em EventMeta) Timestamp() time.Time { return em.ts }
