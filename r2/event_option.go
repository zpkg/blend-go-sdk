package r2

import (
	"net/http"
	"time"

	"github.com/blend/go-sdk/logger"
)

// EventOption is an event option.
type EventOption func(e *Event)

// OptEventFlag sets the event flag.
func OptEventFlag(flag logger.Flag) EventOption {
	return func(e *Event) {
		e.SetFlag(flag)
	}
}

// OptEventCompleted sets the event completed time.
func OptEventCompleted(ts time.Time) EventOption {
	return func(e *Event) {
		e.SetTimestamp(ts)
	}
}

// OptEventStarted sets the start time.
func OptEventStarted(ts time.Time) EventOption {
	return func(e *Event) {
		e.Started = ts
	}
}

// OptEventRequest sets the response.
func OptEventRequest(req *http.Request) EventOption {
	return func(e *Event) {
		e.Request = req
	}
}

// OptEventResponse sets the response.
func OptEventResponse(res *http.Response) EventOption {
	return func(e *Event) {
		e.Response = res
	}
}

// OptEventBody sets the body.
func OptEventBody(body []byte) EventOption {
	return func(e *Event) {
		e.Body = body
	}
}
