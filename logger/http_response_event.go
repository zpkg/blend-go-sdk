package logger

import (
	"bytes"
	"net/http"
	"time"
)

var (
	_ Event = (*HTTPResponseEvent)(nil)
)

// NewHTTPResponseEvent is an event representing a response to an http request.
func NewHTTPResponseEvent(req *http.Request) *HTTPResponseEvent {
	return &HTTPResponseEvent{
		EventMeta: NewEventMeta(HTTPResponse),
		Request:   req,
	}
}

// NewHTTPResponseEventListener returns a new web request event listener.
func NewHTTPResponseEventListener(listener func(*HTTPResponseEvent)) Listener {
	return func(e Event) {
		if typed, isTyped := e.(*HTTPResponseEvent); isTyped {
			listener(typed)
		}
	}
}

// HTTPResponseEvent is an event type for responses.
type HTTPResponseEvent struct {
	*EventMeta

	Request         *http.Request
	Route           string
	ContentLength   int
	ContentType     string
	ContentEncoding string
	StatusCode      int
	Elapsed         time.Duration
	State           map[interface{}]interface{}
}

// WriteText implements TextWritable.
func (e *HTTPResponseEvent) WriteText(formatter TextFormatter, buf *bytes.Buffer) {
	WriteHTTPResponse(formatter, buf, e.Request, e.StatusCode, e.ContentLength, e.ContentType, e.Elapsed)
}

// WriteJSON implements JSONWritable.
func (e *HTTPResponseEvent) WriteJSON() Fields {
	return HTTPResponseFields(e.Request, e.StatusCode, e.ContentLength, e.ContentType, e.ContentEncoding, e.Elapsed)
}
