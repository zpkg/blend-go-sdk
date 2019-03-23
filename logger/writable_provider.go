package logger

// WritableProvider is a provider for if we should write an event to output.
type WritableProvider interface {
	IsWritable() bool
}
