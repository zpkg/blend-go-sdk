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
		Flag:        flag,
		Time:        time.Now().UTC(),
		Labels:      make(map[string]string),
		Annotations: make(map[string]string),
	}
}

// EventMeta is the metadata common to events.
// It is useful for ensuring you have standard fields on your events, and its typically embedded in event types.
type EventMeta struct {
	Flag        string
	FlagColor   ansi.Color
	Time        time.Time
	Labels      map[string]string
	Annotations map[string]string
}
