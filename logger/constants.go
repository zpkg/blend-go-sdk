package logger

import "time"

// Flags
const (
	All          = "all"
	None         = "none"
	Fatal        = "fatal"
	Error        = "error"
	Warning      = "warning"
	Debug        = "debug"
	Info         = "info"
	Silly        = "silly"
	HTTPRequest  = "http.request"
	HTTPResponse = "http.response"
	Audit        = "audit"
	Query        = "db.query"
	RPC          = "rpc"
)

// Default flags
var (
	DefaultFlags = []string{Info, Error, Fatal}
)

// Environment Variable Names
const (
	EnvVarEventFlags       = "LOG_EVENTS"
	EnvVarHiddenEventFlags = "LOG_HIDDEN"
	EnvVarFormat           = "LOG_FORMAT"
	EnvVarUseColor         = "LOG_USE_COLOR"
	EnvVarShowTime         = "LOG_SHOW_TIME"
	EnvVarShowHeadings     = "LOG_SHOW_HEADINGS"
	EnvVarHeading          = "LOG_HEADING"
	EnvVarTimeFormat       = "LOG_TIME_FORMAT"
	EnvVarJSONPretty       = "LOG_JSON_PRETTY"
)

const (
	// Gigabyte is an SI unit.
	Gigabyte int = 1 << 30
	// Megabyte is an SI unit.
	Megabyte int = 1 << 20
	// Kilobyte is an SI unit.
	Kilobyte int = 1 << 10
)

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

const (
	// DefaultWriteQueueDepth  is the default depth per listener to queue work.
	// It's currently set to 256k entries.
	DefaultWriteQueueDepth = 1 << 18

	// DefaultListenerQueueDepth is the default depth per listener to queue work.
	// It's currently set to 256k entries.
	DefaultListenerQueueDepth = 1 << 10
)

// Rune constants
const (
	RuneSpace   rune = ' '
	RuneNewline rune = '\n'
)
