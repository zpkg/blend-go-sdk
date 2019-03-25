package logger

import (
	"io"
	"net/http"
)

// these are compile time assertions
var (
	_ Event          = (*HTTPRequestEvent)(nil)
	_ TextWritable   = (*HTTPRequestEvent)(nil)
	_ FieldsProvider = (*HTTPRequestEvent)(nil)
)

// NewHTTPRequestEvent creates a new web request event.
func NewHTTPRequestEvent(req *http.Request) *HTTPRequestEvent {
	return &HTTPRequestEvent{
		EventMeta: NewEventMeta(HTTPRequest),
		Request:   req,
	}
}

// NewHTTPRequestEventListener returns a new web request event listener.
func NewHTTPRequestEventListener(listener func(*HTTPRequestEvent)) Listener {
	return func(e Event) {
		if typed, isTyped := e.(*HTTPRequestEvent); isTyped {
			listener(typed)
		}
	}
}

// HTTPRequestEvent is an event type for http responses.
type HTTPRequestEvent struct {
	*EventMeta
	Request *http.Request
	Route   string
	State   map[interface{}]interface{}
}

// WriteText implements TextWritable.
func (e *HTTPRequestEvent) WriteText(formatter Colorizer, wr io.Writer) {
	WriteHTTPRequest(formatter, wr, e.Request)
}

// Fields implements the required method to write as json.
func (e *HTTPRequestEvent) Fields() Fields {
	return HTTPRequestFields(e.Request)
}
