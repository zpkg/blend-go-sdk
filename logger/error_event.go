package logger

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/blend/go-sdk/exception"
)

// these are compile time assertions
var (
	_ Event          = (*ErrorEvent)(nil)
	_ TextWritable   = (*ErrorEvent)(nil)
	_ json.Marshaler = (*ErrorEvent)(nil)
)

// Errorf returns a new error event based on format and arguments.
func Errorf(flag, format string, args ...interface{}) *ErrorEvent {
	return &ErrorEvent{
		EventMeta: NewEventMeta(flag),
		Err:       fmt.Errorf(format, args...),
	}
}

// NewErrorEvent returns a new error event.
func NewErrorEvent(flag string, err error, options ...EventMetaOption) *ErrorEvent {
	return &ErrorEvent{
		EventMeta: NewEventMeta(flag, options...),
		Err:       err,
	}
}

// NewErrorEventListener returns a new error event listener.
func NewErrorEventListener(listener func(context.Context, *ErrorEvent)) Listener {
	return func(ctx context.Context, e Event) {
		if typed, isTyped := e.(*ErrorEvent); isTyped {
			listener(ctx, typed)
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
func (e ErrorEvent) WriteText(formatter TextFormatter, output io.Writer) {
	if e.Err != nil {
		if typed, ok := e.Err.(*exception.Ex); ok {
			io.WriteString(output, typed.String())
		} else {
			io.WriteString(output, e.Err.Error())
		}
	}
}

// MarshalJSON implements json.Marshaler.
func (e ErrorEvent) MarshalJSON() ([]byte, error) {
	return json.Marshal(MergeDecomposed(e.EventMeta.Decompose(), map[string]interface{}{
		"err":   e.Err,
		"state": e.State,
	}))
}
