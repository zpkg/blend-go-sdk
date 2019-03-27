package logger

import (
	"context"
	"fmt"
	"io"
)

// these are compile time assertions
var (
	_ Event = &MessageEvent{}
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
	*EventMeta
	Message string
}

// WriteText implements TextWritable.
func (e *MessageEvent) WriteText(formatter TextFormatter, output io.Writer) {
	io.WriteString(output, e.Message)
}

// Fields implements FieldsProvider.
func (e *MessageEvent) Fields() Fields {
	return Fields{
		FieldMessage: e.Message,
	}
}

// String returns the message event body.
func (e *MessageEvent) String() string {
	return e.Message
}
