package logger

import (
	"fmt"
	"io"
	"time"

	"github.com/blend/go-sdk/timeutil"
)

// these are compile time assertions
var (
	_ Event          = (*TimedEvent)(nil)
	_ TextWritable   = (*TimedEvent)(nil)
	_ FieldsProvider = (*TimedEvent)(nil)
)

// Timedf returns a timed message event.
func Timedf(flag string, elapsed time.Duration, format string, args ...interface{}) *TimedEvent {
	return &TimedEvent{
		EventMeta: NewEventMeta(flag),
		Message:   fmt.Sprintf(format, args...),
		Elapsed:   elapsed,
	}
}

// NewTimedEventListener returns a new timed event listener.
func NewTimedEventListener(listener func(e *TimedEvent)) Listener {
	return func(e Event) {
		if typed, isTyped := e.(*TimedEvent); isTyped {
			listener(typed)
		}
	}
}

// TimedEvent is a message event with an elapsed time.
type TimedEvent struct {
	*EventMeta

	Message string
	Elapsed time.Duration
}

// String implements fmt.Stringer
func (e TimedEvent) String() string {
	return fmt.Sprintf("%s (%v)", e.Message, e.Elapsed)
}

// WriteText implements TextWritable.
func (e TimedEvent) WriteText(tf Colorizer, wr io.Writer) {
	io.WriteString(wr, e.String())
}

// Fields implements FieldsProvider.
func (e TimedEvent) Fields() Fields {
	return Fields{
		FieldMessage: e.Message,
		FieldElapsed: timeutil.Milliseconds(e.Elapsed),
	}
}
