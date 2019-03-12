package r2

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/blend/go-sdk/logger"
)

// NewEvent returns a new event.
func NewEvent(flag logger.Flag, options ...EventOption) *Event {
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

	Started  time.Time
	Request  *http.Request
	Response *http.Response
	Body     []byte
}

// WriteText writes the event to a text writer.
func (e *Event) WriteText(tf logger.TextFormatter, buf *bytes.Buffer) {
	if e.Request != nil && e.Response != nil {
		buf.WriteString(fmt.Sprintf("%s %s %s (%v)", e.Request.Method, e.Request.URL.String(), tf.ColorizeStatusCode(e.Response.StatusCode), e.Timestamp().Sub(e.Started)))
	} else if e.Request != nil {
		buf.WriteString(fmt.Sprintf("%s %s", e.Request.Method, e.Request.URL.String()))
	}
	if e.Body != nil {
		buf.WriteRune(logger.RuneNewline)
		buf.Write(e.Body)
	}
}

// WriteJSON implements logger.JSONWritable.
func (e *Event) WriteJSON() logger.JSONObj {
	output := logger.JSONObj{}
	if e.Request != nil {
		output["req"] = logger.JSONObj{
			"startTime": e.Started,
			"method":    e.Request.Method,
			"url":       e.Request.URL.String(),
			"headers":   e.Request.Header,
		}
	}
	if e.Response != nil {
		output["res"] = logger.JSONObj{
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
