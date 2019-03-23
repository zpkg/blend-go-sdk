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
func NewEventMeta(flag string) *EventMeta {
	return &EventMeta{
		EventFlag: flag,
		Time:      time.Now().UTC(),
	}
}

// EventMeta is the metadata common to events.
// It is useful for ensuring you have the minimum required fields on your events, and its typically embedded in types.
type EventMeta struct {
	EventFlag string
	Time      time.Time
	FlagColor ansi.Color
}

// Flag returns the event flag.
func (em EventMeta) Flag() string { return em.EventFlag }

// Timestamp returns the event timestamp.
func (em EventMeta) Timestamp() time.Time { return em.Time }
