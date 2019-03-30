package secrets

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/blend/go-sdk/ansi"
	"github.com/blend/go-sdk/logger"
)

var (
	_ logger.Event        = (*Event)(nil)
	_ logger.TextWritable = (*Event)(nil)
	_ json.Marshaler      = (*Event)(nil)
)

const (
	// Flag is the logger flag.
	Flag = "secrets"
)

// NewEvent returns a new event from a request.
func NewEvent(req *http.Request) *Event {
	return &Event{
		EventMeta: logger.NewEventMeta(Flag),
		Remote:    req.URL.Host,
		Method:    req.Method,
		Key:       strings.TrimPrefix(req.URL.Path, "/v1/"),
	}
}

// Event is an event.
type Event struct {
	*logger.EventMeta
	Remote string
	Method string
	Key    string
}

// MarshalJSON impements json.Marshaler.
func (e *Event) MarshalJSON() ([]byte, error) {
	return json.Marshal(logger.MergeDecomposed(e.EventMeta.Decompose(), map[string]interface{}{
		"remote": e.Remote,
		"method": e.Method,
		"key":    e.Key,
	}))
}

// WriteText writes text for the event.
func (e *Event) WriteText(tf logger.TextFormatter, wr io.Writer) {
	io.WriteString(wr, "["+tf.Colorize(e.Method, ansi.ColorBlue)+"]")
	io.WriteString(wr, logger.Space)
	io.WriteString(wr, e.Remote)
	io.WriteString(wr, logger.Space)
	io.WriteString(wr, e.Key)
}
