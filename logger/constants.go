package logger

import "time"

const (
	// DefaultBufferPoolSize is the default buffer pool size.
	DefaultBufferPoolSize = 1 << 8 // 256

	// DefaultTextTimeFormat is the default time format.
	DefaultTextTimeFormat = time.RFC3339Nano

	// DefaultTextWriterUseColor is a default setting for writers.
	DefaultTextWriterUseColor = true
	// DefaultTextWriterShowHeadings is a default setting for writers.
	DefaultTextWriterShowHeadings = true
	// DefaultTextWriterShowTimestamp is a default setting for writers.
	DefaultTextWriterShowTimestamp = true
)

var (
	// DefaultFlags are the default flags.
	DefaultFlags = []Flag{Fatal, Error, Warning, Info, WebRequest}
	// DefaultFlagSet is the default verbosity for a diagnostics agent inited from the environment.
	DefaultFlagSet = NewFlagSet(DefaultFlags...)

	// DefaultHiddenFlags are the default hidden flags.
	DefaultHiddenFlags []Flag
)

const (
	// DefaultWorkerQueueDepth is the default depth per listener to queue work.
	// It's currently set to 1 million entries.
	DefaultWorkerQueueDepth = 1 << 20
)
