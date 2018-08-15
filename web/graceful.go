package web

import (
	"os"
	"os/signal"
	"syscall"
)

// GracefulShutdown starts an app and responds to SIGINT and SIGTERM to shut the app down.
func GracefulShutdown(app *App) {
	shutdown := make(chan struct{})
	serverExit := make(chan struct{})
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		if err := app.Start(); err != nil {
			if app.log != nil {
				app.log.SyncFatal(err)
			}
		}
		close(serverExit)
	}()

	go func() {
		<-quit
		if err := app.Shutdown(); err != nil {
			app.log.SyncFatal(err)
		}
		close(shutdown)
	}()

	<-shutdown
	<-serverExit
}
