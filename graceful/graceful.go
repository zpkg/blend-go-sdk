package graceful

import (
	"os"
	"os/signal"
	"sync"
	"syscall"
)

// Graceful is a server that can start and shutdown.
type Graceful interface {
	// Start the service. This must block.
	Start() error
	// Stop the service.
	Stop() error
	// Notify the service has started.
	NotifyStarted() <-chan struct{}
	// Notify the service has stopped.
	NotifyStopped() <-chan struct{}
}

// Shutdown starts an hosted process and responds to SIGINT and SIGTERM to shut the app down.
// It will return any errors returned by app.Start() that are not caused by shutting down the server.
func Shutdown(hosted ...Graceful) error {
	terminateSignal := make(chan os.Signal, 1)
	signal.Notify(terminateSignal, os.Interrupt, syscall.SIGTERM)
	return ShutdownBySignal(terminateSignal, hosted...)
}

// ShutdownBySignal gracefully stops a hosted process based on an os signal channel.
// A "Graceful" process *must* block on start.
func ShutdownBySignal(shouldShutdown chan os.Signal, hosted ...Graceful) error {
	shutdown := make(chan struct{})
	abortWaitShutdown := make(chan struct{})
	serverExited := make(chan struct{})

	waitShutdownComplete := sync.WaitGroup{}
	waitShutdownComplete.Add(len(hosted))

	waitServerExited := sync.WaitGroup{}
	waitServerExited.Add(len(hosted))

	errors := make(chan error, 2*len(hosted))

	for _, hostedInstance := range hosted {
		go func(instance Graceful) {
			// signal hosted has exited
			defer close(serverExited)

			// `hosted.Start()` should block here.
			if err := instance.Start(); err != nil {
				errors <- err
			}
		}(hostedInstance)

		go func(instance Graceful) {
			select {
			case <-shutdown:
				// tell the hosted process to terminate "gracefully"
				if err := instance.Stop(); err != nil {
					errors <- err
				}
				waitShutdownComplete.Done()
				return
			case <-abortWaitShutdown:
				waitShutdownComplete.Done()
				return
			}
		}(hostedInstance)
	}

	select {
	case <-shouldShutdown: // if we've issued a shutdown, wait for the server to exit
		close(shutdown)
		waitShutdownComplete.Wait()
		waitServerExited.Wait()
	case <-serverExited: // if any of the servers exited on their own
		close(abortWaitShutdown) // quit the signal listener
		waitShutdownComplete.Wait()
	}
	if len(errors) > 0 {
		return <-errors
	}
	return nil
}
