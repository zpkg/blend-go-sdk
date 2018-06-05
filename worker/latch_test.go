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

	l.Starting()
	assert.True(l.IsStarting())
	assert.False(l.IsRunning())
	assert.False(l.IsStopping())
	assert.False(l.IsStopped())
	go func() {
		l.Started()
		didStart = true
		for {
			select {
			case <-work:
				didGetWork = true
				workComplete <- true
			case <-l.NotifyStop():
				didAbort = true
				l.Stopped()
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
	<-l.NotifyStopped()

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
