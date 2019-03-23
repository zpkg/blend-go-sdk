package logger

import "time"

// Event is an interface representing methods necessary to trigger listeners.
type Event interface {
	Flag() string
	Timestamp() time.Time
}

// MarshalEvent marshals an object as a logger event.
func MarshalEvent(obj interface{}) (Event, bool) {
	typed, isTyped := obj.(Event)
	return typed, isTyped
}
