package worker

import (
	"sync"
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
)

func TestIntervalWorker(t *testing.T) {
	assert := assert.New(t)

	var didWork bool
	wg := sync.WaitGroup{}
	wg.Add(1)
	w := NewInterval(func() error {
		didWork = true
		wg.Done()
		return nil
	}, time.Millisecond)

	w.Start()
	assert.True(w.Latch().IsRunning())
	wg.Wait()
	w.Stop()
	assert.True(w.Latch().IsStopped())

	assert.True(didWork)
}
