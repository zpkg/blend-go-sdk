package r2

import (
	"bytes"
	"net/http"
	"time"

	"github.com/blend/go-sdk/logger"
)

// EventOption is an event option.
type EventOption func(e *Event)

// EventFlag sets the event flag.
func EventFlag(flag logger.Flag) EventOption {
	return func(e *Event) {
		e.SetFlag(flag)
	}
}

// EventCompleted sets the event completed time.
func EventCompleted(ts time.Time) EventOption {
	return func(e *Event) {
		e.SetTimestamp(ts)
	}
}

// EventStarted sets the start time.
func EventStarted(ts time.Time) EventOption {
	return func(e *Event) {
		e.Started = ts
	}
}

// EventRequest sets the response.
func EventRequest(req *http.Request) EventOption {
	return func(e *Event) {
		e.Request = req
	}
}

// EventResponse sets the response.
func EventResponse(res *http.Response) EventOption {
	return func(e *Event) {
		e.Response = res
	}
}

// EventBody sets the body.
func EventBody(body *bytes.Buffer) EventOption {
	return func(e *Event) {
		e.Body = body
	}
}
