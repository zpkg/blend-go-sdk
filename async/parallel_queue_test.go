package async

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestParallelQueue(t *testing.T) {
	t.Skip()
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
	errors := make(chan error, 8)
	w := NewParallelQueue(4, func(obj interface{}) error {
		atomic.AddInt32(&finished, 1)
		return fmt.Errorf("only a test %d", finished)
	}).WithErrors(errors)
	w.Start()
	assert.True(w.Latch().IsRunning())

	for x := 0; x < 8; x++ {
		w.Enqueue("hello")
	}
	w.Drain()
	assert.Equal(8, finished)
	assert.Equal(8, len(w.Errors()))
}
