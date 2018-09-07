package web

import (
	"os"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestGracefulShutdown(t *testing.T) {
	assert := assert.New(t)

	app := New()
	terminateSignal := make(chan os.Signal)
	var err error
	done := make(chan struct{})
	go func() {
		err = startWithGracefulShutdownBySignal(app, terminateSignal)
		close(done)
	}()
	<-app.NotifyStarted()

	close(terminateSignal)
	<-done
	assert.Nil(err)
}
