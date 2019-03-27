package logger

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/blend/go-sdk/webutil"
)

// these are compile time assertions
var (
	_ Event          = (*HTTPRequestEvent)(nil)
	_ TextWritable   = (*HTTPRequestEvent)(nil)
	_ json.Marshaler = (*HTTPRequestEvent)(nil)
)

// NewHTTPRequestEvent creates a new web request event.
func NewHTTPRequestEvent(req *http.Request) *HTTPRequestEvent {
	return &HTTPRequestEvent{
		EventMeta: NewEventMeta(HTTPRequest),
		Request:   req,
	}
}

// NewHTTPRequestEventListener returns a new web request event listener.
func NewHTTPRequestEventListener(listener func(context.Context, *HTTPRequestEvent)) Listener {
	return func(ctx context.Context, e Event) {
		if typed, isTyped := e.(*HTTPRequestEvent); isTyped {
			listener(ctx, typed)
		}
	}
}

// HTTPRequestEvent is an event type for http responses.
type HTTPRequestEvent struct {
	*EventMeta `json:",inline"`
	Request    *http.Request
	Route      string
	State      map[interface{}]interface{}
}

// WriteText implements TextWritable.
func (e *HTTPRequestEvent) WriteText(formatter TextFormatter, wr io.Writer) {
	WriteHTTPRequest(formatter, wr, e.Request)
}

// MarshalJSON marshals the event as json.
func (e *HTTPRequestEvent) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"verb":      e.Request.Method,
		"path":      e.Request.URL.Path,
		"host":      e.Request.Host,
		"ip":        webutil.GetRemoteAddr(e.Request),
		"userAgent": webutil.GetUserAgent(e.Request),
	})
}
