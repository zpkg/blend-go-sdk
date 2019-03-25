package logger

import "io"

// TextWritable is an event that can be written.
type TextWritable interface {
	WriteText(Colorizer, io.Writer)
}
