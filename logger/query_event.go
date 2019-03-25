package logger

import (
	"fmt"
	"io"
	"time"

	"github.com/blend/go-sdk/stringutil"

	"github.com/blend/go-sdk/ansi"
	"github.com/blend/go-sdk/timeutil"
)

// these are compile time assertions
var (
	_ Event = (*QueryEvent)(nil)
)

// NewQueryEvent creates a new query event.
func NewQueryEvent(body string, elapsed time.Duration) *QueryEvent {
	return &QueryEvent{
		EventMeta: NewEventMeta(Query),
		Body:      body,
		Elapsed:   elapsed,
	}
}

// NewQueryEventListener returns a new listener for spiffy events.
func NewQueryEventListener(listener func(e *QueryEvent)) Listener {
	return func(e Event) {
		if typed, isTyped := e.(*QueryEvent); isTyped {
			listener(typed)
		}
	}
}

// QueryEvent represents a database query.
type QueryEvent struct {
	*EventMeta

	Database   string
	Engine     string
	Username   string
	QueryLabel string
	Body       string
	Elapsed    time.Duration
	Err        error
}

// WriteText writes the event text to the output.
func (e QueryEvent) WriteText(tf TextFormatter, wr io.Writer) {
	io.WriteString(wr, "[")
	if len(e.Engine) > 0 {
		io.WriteString(wr, tf.Colorize(e.Engine, ansi.ColorLightWhite))
		io.WriteString(wr, Space)
	}
	if len(e.Username) > 0 {
		io.WriteString(wr, tf.Colorize(e.Username, ansi.ColorLightWhite))
		io.WriteString(wr, "@")
	}
	io.WriteString(wr, tf.Colorize(e.Database, ansi.ColorLightWhite))
	io.WriteString(wr, "]")

	if len(e.QueryLabel) > 0 {
		io.WriteString(wr, Space)
		io.WriteString(wr, fmt.Sprintf("[%s]", tf.Colorize(e.QueryLabel, ansi.ColorLightWhite)))
	}

	io.WriteString(wr, Space)
	io.WriteString(wr, e.Elapsed.String())

	if e.Err != nil {
		io.WriteString(wr, Space)
		io.WriteString(wr, tf.Colorize("failed", ansi.ColorRed))
	}

	if len(e.Body) > 0 {
		io.WriteString(wr, Space)
		io.WriteString(wr, stringutil.CompressSpace(e.Body))
	}
}

// Fields implements FieldsProvider.
func (e QueryEvent) Fields() Fields {
	return Fields{
		"engine":     e.Engine,
		"database":   e.Database,
		"username":   e.Username,
		"queryLabel": e.QueryLabel,
		"body":       e.Body,
		FieldErr:     e.Err,
		FieldElapsed: timeutil.Milliseconds(e.Elapsed),
	}
}
