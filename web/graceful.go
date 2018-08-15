package web

import (
	"os"
	"os/signal"
	"syscall"
)

// GracefulShutdown starts an app and responds to SIGINT and SIGTERM to shut the app down.
// It will return any errors returned by app.Start() that are not caused by shutting down the server.
func GracefulShutdown(app *App) error {
	shutdown := make(chan struct{})
	server := make(chan struct{})
	abort := make(chan struct{})

	terminate := make(chan os.Signal, 1)
	signal.Notify(terminate, os.Interrupt, syscall.SIGTERM)

	errors := make(chan error, 2)

	go func() {
		if err := app.Start(); err != nil {
			errors <- err
		}
		close(server)
	}()

	go func() {
		select {
		case <-terminate:
			if err := app.Shutdown(); err != nil {
				errors <- err
			}
			close(shutdown)
		case <-abort:
			return
		}
	}()

	select {
	case <-shutdown: // if we've issued a shutdown, wait for the server to exit
		<-server
	case <-server: // if the server exited
		close(abort) // quit the signal listener
	}

	if len(errors) > 0 {
		return <-errors
	}
	return nil
}
