package async

import (
	"sync"
	"sync/atomic"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestParallelQueue(t *testing.T) {
	assert := assert.New(t)

	wg := sync.WaitGroup{}
	wg.Add(8)
	w := NewParallelQueue(4, func(obj interface{}) error {
		defer wg.Done()
		return nil
	})
	w.Start()
	assert.True(w.Latch().IsRunning())

	for x := 0; x < 8; x++ {
		w.Enqueue("hello")
	}

	wg.Wait()
	w.Close()
	assert.False(w.Latch().IsRunning())
}

func TestParallelQueueDrain(t *testing.T) {
	assert := assert.New(t)

	var finished int32
	w := NewParallelQueue(4, func(obj interface{}) error {
		atomic.AddInt32(&finished, 1)
		return nil
	})
	w.Start()
	assert.True(w.Latch().IsRunning())

	for x := 0; x < 8; x++ {
		w.Enqueue("hello")
	}
	w.Drain()
	assert.Equal(8, finished)
}
