package timeutil

import "time"

// Milliseconds returns a duration as milliseconds.
func Milliseconds(d time.Duration) float64 {
	return float64(d) / float64(time.Millisecond)
}

// FromMilliseconds returns a duration from a given float64 millis value.
func FromMilliseconds(millis float64) time.Duration {
	return time.Duration(millis * float64(time.Millisecond))
}
