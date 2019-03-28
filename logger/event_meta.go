package logger

import (
	"fmt"
	"time"

	"github.com/blend/go-sdk/ansi"
)

// these are compile time assertions
var (
	_ Event          = (*EventMeta)(nil)
	_ FieldsProvider = (*EventMeta)(nil)
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

// OptEventMetaFlagColor sets the event flag color.
func OptEventMetaFlagColor(color ansi.Color) EventMetaOption {
	return func(em *EventMeta) { em.flagColor = color }
}

// OptEventMetaFlag sets the event flag.
func OptEventMetaFlag(flag string) EventMetaOption {
	return func(em *EventMeta) { em.flag = flag }
}

// OptEventMetaTimestamp sets the event timestamp.
func OptEventMetaTimestamp(ts time.Time) EventMetaOption {
	return func(em *EventMeta) { em.ts = ts }
}

// OptEventMetaField sets an event meta field.
func OptEventMetaField(key string, value interface{}) EventMetaOption {
	return func(em *EventMeta) {
		if em.fields == nil {
			em.fields = make(map[string]string)
		}
		em.fields[key] = fmt.Sprintf("%v", value)
	}
}

// EventMeta is the metadata common to events.
// It is useful for ensuring you have the minimum required fields on your events, and its typically embedded in types.
type EventMeta struct {
	flag      string
	ts        time.Time
	flagColor ansi.Color
	fields    map[string]string
}

// Flag returns the event flag.
func (em EventMeta) Flag() string { return em.flag }

// FlagColor returns the event flag color
func (em EventMeta) FlagColor() ansi.Color { return em.flagColor }

// Timestamp returns the event timestamp.
func (em EventMeta) Timestamp() time.Time { return em.ts }

// Fields returns the event meta fields.
func (em EventMeta) Fields() map[string]string { return em.fields }

// Decompose decomposes the object into a map[string]interface{}.
func (em EventMeta) Decompose() map[string]interface{} {
	output := map[string]interface{}{
		FieldFlag:      em.flag,
		FieldTimestamp: em.ts.Format(time.RFC3339Nano),
	}
	if em.fields != nil {
		output[FieldFields] = em.fields
	}
	return output
}
