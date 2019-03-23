package logger

import (
	"encoding/json"
	"io"
)

// NewJSONFormatter returns a new json event formatter.
func NewJSONFormatter(cfg *JSONConfig) *JSONFormatter {
	return &JSONFormatter{
		Pretty: cfg.Pretty,
	}
}

// JSONFormatter is a json output formatter.
type JSONFormatter struct {
	Pretty  bool
	Encoder *json.Encoder
}

// WriteFormat writes the event to the given output.
func (jw JSONFormatter) WriteFormat(output io.Writer, e Event) error {
	encoder := json.NewEncoder(output)
	if jw.Pretty {
		encoder.SetIndent("", "\t")
	}

	if typed, isTyped := e.(FieldsProvider); isTyped {
		fields := typed.Fields()
		fields[FieldFlag] = e.Flag()
		fields[FieldTimestamp] = e.Timestamp()
		return encoder.Encode(fields)
	}

	return encoder.Encode(e)
}
