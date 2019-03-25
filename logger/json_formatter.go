package logger

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
)

// NewJSONFormatter returns a new json event formatter.
func NewJSONFormatter(options ...JSONFormatterOption) *JSONFormatter {
	jf := &JSONFormatter{}

	for _, option := range options {
		option(jf)
	}
	return jf
}

// JSONFormatterOption is an option for json formatters.
type JSONFormatterOption func(*JSONFormatter)

// OptJSONConfig sets a json formatter from a config.
func OptJSONConfig(cfg *JSONConfig) JSONFormatterOption {
	return func(jf *JSONFormatter) {
		jf.Pretty = cfg.Pretty
		jf.PrettyIndent = cfg.PrettyIndentOrDefault()
		jf.PrettyPrefix = cfg.PrettyPrefixOrDefault()
	}
}

// JSONFormatter is a json output formatter.
type JSONFormatter struct {
	Pretty       bool
	PrettyPrefix string
	PrettyIndent string
}

// WriteFormat writes the event to the given output.
func (jw JSONFormatter) WriteFormat(ctx context.Context, output io.Writer, e Event) error {
	buffer := new(bytes.Buffer)
	encoder := json.NewEncoder(buffer)
	if jw.Pretty {
		encoder.SetIndent(jw.PrettyPrefix, jw.PrettyIndent)
	}

	if typed, isTyped := e.(FieldsProvider); isTyped {
		fields := typed.Fields()
		fields[FieldFlag] = e.Flag()
		fields[FieldTimestamp] = e.Timestamp()
		if err := encoder.Encode(fields); err != nil {
			return err
		}
	}

	_, err := io.Copy(output, buffer)
	return err
}
