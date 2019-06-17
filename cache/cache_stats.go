package cache

import "time"

// Stats represents cached statistics.
type Stats struct {
	Count     int
	SizeBytes int
	MaxAge    time.Duration
}
