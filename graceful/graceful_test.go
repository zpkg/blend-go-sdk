package graceful

import (
	"os"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func newHosted() *hosted {
	return &hosted{
		started: make(chan struct{}),
	}
}

type hosted struct {
	state   int32
	started chan struct{}
	stopped chan struct{}
}

func (h *hosted) Start() error {
	h.state = 1
	h.stopped = make(chan struct{})
	close(h.started)
	return nil
}

func (h *hosted) Stop() error {
	h.state = 0
	h.started = make(chan struct{})
	close(h.stopped)
	return nil
}

func (h *hosted) NotifyStarted() <-chan struct{} {
	return h.started
}

func (h *hosted) NotifyStopped() <-chan struct{} {
	return h.stopped
}

func TestGracefulShutdown(t *testing.T) {
	assert := assert.New(t)

	hosted := newHosted()
	terminateSignal := make(chan os.Signal)
	var err error
	done := make(chan struct{})
	go func() {
		err = ShutdownBySignal(hosted, terminateSignal)
		close(done)
	}()
	<-hosted.NotifyStarted()

	close(terminateSignal)
	<-done
	assert.Nil(err)
}
