package secrets

import (
	"bytes"
	"net/http"
	"strings"

	"github.com/blend/go-sdk/logger"
)

const (
	// Flag is the logger flag.
	Flag = "secret"
)

// NewEvent returns a new event from a request.
func NewEvent(req *http.Request) *Event {
	return &Event{
		EventMeta: logger.NewEventMeta(Flag),
		remote:    req.URL.Host,
		method:    req.Method,
		key:       strings.TrimPrefix(req.URL.Path, "/v1/"),
	}
}

// Event is an event.
type Event struct {
	*logger.EventMeta

	remote string
	method string
	key    string
}

// WithRemote sets the event remote.
func (e *Event) WithRemote(remote string) *Event {
	e.remote = remote
	return e
}

// Remote returns the remote.
func (e *Event) Remote() string {
	return e.remote
}

// WithMethod sets the event method.
func (e *Event) WithMethod(method string) *Event {
	e.method = method
	return e
}

// Method returns the method.
func (e *Event) Method() string {
	return e.method
}

// WithKey sets the event method.
func (e *Event) WithKey(key string) *Event {
	e.key = key
	return e
}

// Key returns the event key.
func (e *Event) Key() string {
	return e.key
}

// WriteJSON returns json values.
func (e *Event) WriteJSON() map[string]interface{} {
	return map[string]interface{}{
		"method": e.method,
		"key":    e.key,
	}
}

// WriteText writes text for the event.
func (e *Event) WriteText(tf logger.TextFormatter, buf *bytes.Buffer) {
	buf.WriteString("[" + tf.Colorize(e.method, logger.ColorBlue) + "]")
	buf.WriteRune(logger.RuneSpace)
	buf.WriteString(e.remote)
	buf.WriteRune(logger.RuneSpace)
	buf.WriteString(e.key)
}
