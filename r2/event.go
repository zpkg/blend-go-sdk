package r2

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/blend/go-sdk/logger"
)

const (
	// Flag is a logger event flag.
	Flag = "request"
	// FlagResponse is a logger event flag.
	FlagResponse = "request.response"
)

// NewEvent returns a new event.
func NewEvent(flag string, options ...EventOption) *Event {
	e := &Event{
		EventMeta: logger.NewEventMeta(flag),
	}
	for _, option := range options {
		option(e)
	}
	return e
}

// Event is a response to outgoing requests.
type Event struct {
	*logger.EventMeta

	// Started is the time the request was started.
	// It is used for elapsed time calculations.
	Started time.Time
	// The request metadata.
	Request *http.Request
	// The response metadata (excluding the body).
	Response *http.Response
	// The response body.
	Body []byte
}

// WriteText writes the event to a text writer.
func (e *Event) WriteText(tf logger.Colorizer, wr io.Writer) {
	if e.Request != nil && e.Response != nil {
		io.WriteString(wr, fmt.Sprintf("%s %s %s (%v)", e.Request.Method, e.Request.URL.String(), logger.ColorizeStatusCode(tf, e.Response.StatusCode), e.Timestamp().Sub(e.Started)))
	} else if e.Request != nil {
		io.WriteString(wr, fmt.Sprintf("%s %s", e.Request.Method, e.Request.URL.String()))
	}
	if e.Body != nil {
		io.WriteString(wr, logger.Newline)
		io.WriteString(wr, string(e.Body))
	}
}

// Fields implements logger.FieldsProvider.
func (e *Event) Fields() logger.Fields {
	output := make(logger.Fields)
	if e.Request != nil {
		output["req"] = logger.Fields{
			"startTime": e.Started,
			"method":    e.Request.Method,
			"url":       e.Request.URL.String(),
			"headers":   e.Request.Header,
		}
	}
	if e.Response != nil {
		output["res"] = logger.Fields{
			"completeTime":    e.Timestamp(),
			"statusCode":      e.Response.StatusCode,
			"contentLength":   e.Response.ContentLength,
			"contentType":     tryHeader(e.Response.Header, "Content-Type", "content-type"),
			"contentEncoding": tryHeader(e.Response.Header, "Content-Encoding", "content-encoding"),
			"headers":         e.Response.Header,
			"cert":            ParseCertInfo(e.Response),
		}
	}
	if e.Body != nil {
		output["body"] = string(e.Body)
	}

	return output
}

func tryHeader(headers http.Header, keys ...string) string {
	for _, key := range keys {
		if values, hasValues := headers[key]; hasValues {
			return strings.Join(values, ";")
		}
	}
	return ""
}
