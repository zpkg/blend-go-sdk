package timeutil

import "time"

// Milliseconds returns a duration as milliseconds.
func Milliseconds(d time.Duration) float64 {
	return float64(d) / float64(time.Millisecond)
}
