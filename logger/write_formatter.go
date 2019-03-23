package logger

import "io"

// WriteFormatter is a formatter for writing events to output writers.
type WriteFormatter interface {
	WriteFormat(io.Writer, Event) error
}
