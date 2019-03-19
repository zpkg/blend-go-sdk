package async

import "time"

// Latch states
const (
	LatchStopped  int32 = 0
	LatchStarting int32 = 1
	LatchRunning  int32 = 2
	LatchStopping int32 = 3
)

// Constants
const (
	DefaultQueueMaxWork = 1 << 10
	DefaultInterval     = 500 * time.Millisecond
)
