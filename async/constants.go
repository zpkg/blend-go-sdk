package async

import "time"

// Latch states
const (
	LatchStopped  int32 = 0
	LatchStarting int32 = 1
	LatchStarted  int32 = 2
	LatchPausing  int32 = 3
	LatchPaused   int32 = 4
	LatchResuming int32 = 5
	LatchStopping int32 = 6
)

// Constants
const (
	DefaultQueueMaxWork = 1 << 10
	DefaultInterval     = 500 * time.Millisecond
)
