package logger

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
)

// these are compile time assertions
var (
	_ Event = (*MessageEvent)(nil)
)

// Messagef returns a new Message Event.
func Messagef(flag, format string, args ...interface{}) *MessageEvent {
	return &MessageEvent{
		EventMeta: NewEventMeta(flag),
		Message:   fmt.Sprintf(format, args...),
	}
}

// NewMessageEventListener returns a new message event listener.
func NewMessageEventListener(listener func(context.Context, *MessageEvent)) Listener {
	return func(ctx context.Context, e Event) {
		if typed, isTyped := e.(*MessageEvent); isTyped {
			listener(ctx, typed)
		}
	}
}

// MessageEvent is a common type of message.
type MessageEvent struct {
	*EventMeta `json:",inline"`
	Message    string `json:"message"`
}

// WriteText implements TextWritable.
func (e *MessageEvent) WriteText(formatter TextFormatter, output io.Writer) {
	io.WriteString(output, e.Message)
}

// MarshalJSON implements json.Marshaler.
func (e MessageEvent) MarshalJSON() ([]byte, error) {
	return json.Marshal(MergeDecomposed(e.EventMeta.Decompose(), map[string]interface{}{
		"message": e.Message,
	}))
}
