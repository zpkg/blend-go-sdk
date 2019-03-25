package logger

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/blend/go-sdk/exception"
)

// these are compile time assertions
var (
	_ Event          = (*ErrorEvent)(nil)
	_ TextWritable   = (*ErrorEvent)(nil)
	_ FieldsProvider = (*ErrorEvent)(nil)
)

// Errorf returns a new error event based on format and arguments.
func Errorf(flag, format string, args ...interface{}) *ErrorEvent {
	return &ErrorEvent{
		EventMeta: NewEventMeta(flag),
		Err:       fmt.Errorf(format, args...),
	}
}

// NewErrorEvent returns a new error event.
func NewErrorEvent(flag string, err error) *ErrorEvent {
	return &ErrorEvent{
		EventMeta: NewEventMeta(flag),
		Err:       err,
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

// WriteText writes the text version of an error.
func (e *ErrorEvent) WriteText(formatter Colorizer, output io.Writer) {
	if e.Err != nil {
		if typed, ok := e.Err.(*exception.Ex); ok {
		} else {
			io.WriteString(output, e.Err.Error())
		}
	}
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
