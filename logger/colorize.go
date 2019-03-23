package logger

import (
	"net/http"
	"strconv"

	"github.com/blend/go-sdk/ansi"
)

var (
	// DefaultFlagTextColors is the default color for each known flag.
	DefaultFlagTextColors = map[string]ansi.Color{
		Info:    ansi.ColorLightWhite,
		Silly:   ansi.ColorLightBlack,
		Debug:   ansi.ColorLightYellow,
		Warning: ansi.ColorLightYellow,
		Error:   ansi.ColorRed,
		Fatal:   ansi.ColorRed,
	}

	// DefaultFlagTextColor is the default flag color.
	DefaultFlagTextColor = ansi.ColorLightWhite
)

// FlagTextColor returns the color for a flag.
func FlagTextColor(flag string) ansi.Color {
	if color, hasColor := DefaultFlagTextColors[flag]; hasColor {
		return color
	}
	return DefaultFlagTextColor
}

// ColorizeByStatusCode returns a value colored by an http status code.
func ColorizeByStatusCode(statusCode int, value string) string {
	if statusCode >= http.StatusOK && statusCode < 300 { //the http 2xx range is ok
		return ansi.ColorGreen.Apply(value)
	} else if statusCode == http.StatusInternalServerError {
		return ansi.ColorRed.Apply(value)
	}
	return ansi.ColorYellow.Apply(value)
}

// ColorizeStatusCode colorizes a status code.
func ColorizeStatusCode(statusCode int) string {
	return ColorizeByStatusCode(statusCode, strconv.Itoa(statusCode))
}
