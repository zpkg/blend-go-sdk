package async

import (
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
)

func TestIntervalWorker(t *testing.T) {
	assert := assert.New(t)

	var didWork bool
	done := make(chan struct{})
	w := NewInterval(func() error {
		defer func() {
			close(done)
		}()
		didWork = true
		return nil
	}, time.Millisecond)

	w.Start()
	assert.True(w.Latch().IsRunning())
	<-done
	w.Stop()
	assert.True(w.Latch().IsStopped())
	assert.True(didWork)
}
