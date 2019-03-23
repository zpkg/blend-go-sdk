package logger

import "github.com/blend/go-sdk/ansi"

// Colorizer is a type that can colorize a given value.
type Colorizer interface {
	Colorize(string, ansi.Color)
}
