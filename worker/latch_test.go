package worker

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestLatch(t *testing.T) {
	assert := assert.New(t)

	l := &Latch{}

	var didStart bool
	var didAbort bool
	var didGetWork bool

	work := make(chan bool)
	workComplete := make(chan bool)

	l.SignalStarting()
	assert.True(l.IsStarting())
	assert.False(l.IsRunning())
	assert.False(l.IsStopping())
	assert.False(l.IsStopped())
	go func() {
		l.SignalStarted()
		didStart = true
		for {
			select {
			case <-work:
				didGetWork = true
				workComplete <- true
			case <-l.ShouldStop():
				didAbort = true
				l.SignalStopped()
				return
			}
		}
	}()

	work <- true
	assert.True(l.IsRunning())

	// wait for work to happen.
	<-workComplete

	// signal stop
	l.Stop()
	<-l.Stopped()

	assert.True(didStart)
	assert.True(didAbort)
	assert.True(didGetWork)
	assert.False(l.IsStopping())
	assert.False(l.IsRunning())
	assert.True(l.IsStopped())

	didAbort = false
	l.Stop()

	assert.False(didAbort)
	assert.False(l.IsStopping())
	assert.False(l.IsRunning())
	assert.True(l.IsStopped())
}
