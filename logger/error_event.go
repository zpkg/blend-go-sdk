package logger

import (
	"encoding/json"
	"fmt"
)

// these are compile time assertions
var (
	_ Event = &ErrorEvent{}
)

// Errorf returns a new error event based on format and arguments.
func Errorf(flag, format string, args ...interface{}) *ErrorEvent {
	return &ErrorEvent{
		EventMeta: NewEventMeta(flag),
		Err:       fmt.Errorf(format, args...),
	}
}

// NewErrorEventListener returns a new error event listener.
func NewErrorEventListener(listener func(*ErrorEvent)) Listener {
	return func(e Event) {
		if typed, isTyped := e.(*ErrorEvent); isTyped {
			listener(typed)
		}
	}
}

// ErrorEvent is an event that wraps an error.
type ErrorEvent struct {
	*EventMeta

	Err   error
	State interface{}
}

// Fields implements FieldsProvider.
func (e *ErrorEvent) Fields() Fields {
	var err interface{}
	if _, ok := e.Err.(json.Marshaler); ok {
		err = e.Err
	} else {
		err = e.Err.Error()
	}
	return Fields{
		FieldErr: err,
	}
}
