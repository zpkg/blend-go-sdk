package web

// Shutdowner is a server that can start and shutdown.
type Shutdowner interface {
	Start() error
	Shutdown() error
	IsRunning() bool
	NotifyStarted() <-chan struct{}
	NotifyShutdown() <-chan struct{}
}
