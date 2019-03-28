package logger

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/blend/go-sdk/timeutil"
	"github.com/blend/go-sdk/webutil"
)

var (
	_ Event          = (*HTTPResponseEvent)(nil)
	_ TextWritable   = (*HTTPResponseEvent)(nil)
	_ json.Marshaler = (*HTTPResponseEvent)(nil)
)

// NewHTTPResponseEvent is an event representing a response to an http request.
func NewHTTPResponseEvent(req *http.Request, options ...EventMetaOption) *HTTPResponseEvent {
	return &HTTPResponseEvent{
		EventMeta: NewEventMeta(HTTPResponse, options...),
		Request:   req,
	}
}

// NewHTTPResponseEventListener returns a new web request event listener.
func NewHTTPResponseEventListener(listener func(context.Context, *HTTPResponseEvent)) Listener {
	return func(ctx context.Context, e Event) {
		if typed, isTyped := e.(*HTTPResponseEvent); isTyped {
			listener(ctx, typed)
		}
	}
}

// HTTPResponseEvent is an event type for responses.
type HTTPResponseEvent struct {
	*EventMeta `json:",inline"`

	Request         *http.Request               `json:"-"`
	Route           string                      `json:"route"`
	ContentLength   int                         `json:"contentLength"`
	ContentType     string                      `json:"contentType"`
	ContentEncoding string                      `json:"contentEncoding"`
	StatusCode      int                         `json:"statusCode"`
	Elapsed         time.Duration               `json:"elapsed"`
	State           map[interface{}]interface{} `json:"state,omitempty"`
}

// WriteText implements TextWritable.
func (e HTTPResponseEvent) WriteText(formatter TextFormatter, wr io.Writer) {
	WriteHTTPResponse(formatter, wr, e.Request, e.StatusCode, e.ContentLength, e.ContentType, e.Elapsed)
}

// MarshalJSON implements json.Marshaler.
func (e HTTPResponseEvent) MarshalJSON() ([]byte, error) {
	return json.Marshal(MergeDecomposed(e.EventMeta.Decompose(), map[string]interface{}{
		"ip":              webutil.GetRemoteAddr(e.Request),
		"userAgent":       webutil.GetUserAgent(e.Request),
		"verb":            e.Request.Method,
		"path":            e.Request.URL.Path,
		"host":            e.Request.Host,
		"contentLength":   e.ContentLength,
		"contentType":     e.ContentType,
		"contentEncoding": e.ContentEncoding,
		"statusCode":      e.StatusCode,
		"elapsed":         timeutil.Milliseconds(e.Elapsed),
	}))
}
