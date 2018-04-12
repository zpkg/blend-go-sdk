package logger

import (
	"time"
)

// Event is an interface representing methods necessary to trigger listeners.
type Event interface {
	Flag() Flag
	Timestamp() time.Time
}

// EventHeading determines if we should add another output field, `event-heading` to output.
type EventHeading interface {
	Heading() string
}

// EventMeta determines if we should pull extra meta fields off the event.
type EventMeta interface {
	Labels() map[string]string
	Annotations() map[string]string
}

// EventEnabled determines if we should allow an event to be triggered or not.
type EventEnabled interface {
	IsEnabled() bool
}

// EventWritable lets us disable implicit writing for some events.
type EventWritable interface {
	IsWritable() bool
}

// EventError determines if we should write the event to the error stream.
type EventError interface {
	IsError() bool
}

// --------------------------------------------------------------------------------
// testing helpers
// --------------------------------------------------------------------------------

// MarshalEvent marshals an object as a logger event.
func MarshalEvent(obj interface{}) (Event, bool) {
	typed, isTyped := obj.(Event)
	return typed, isTyped
}

// MarshalEventHeading marshals an object as an event heading provider.
func MarshalEventHeading(obj interface{}) (EventHeading, bool) {
	typed, isTyped := obj.(EventHeading)
	return typed, isTyped
}

// MarshalEventEnabled marshals an object as an event enabled provider.
func MarshalEventEnabled(obj interface{}) (EventEnabled, bool) {
	typed, isTyped := obj.(EventEnabled)
	return typed, isTyped
}

// MarshalEventWritable marshals an object as an event writable provider.
func MarshalEventWritable(obj interface{}) (EventWritable, bool) {
	typed, isTyped := obj.(EventWritable)
	return typed, isTyped
}

// MarshalEventMeta marshals an object as an event meta provider.
func MarshalEventMeta(obj interface{}) (EventMeta, bool) {
	typed, isTyped := obj.(EventMeta)
	return typed, isTyped
}
