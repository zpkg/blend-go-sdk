package logger

// EnabledProvider is an enabled provider.
type EnabledProvider interface {
	IsEnabled() bool
}
