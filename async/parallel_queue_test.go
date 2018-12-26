package async

import (
	"context"
	"sync"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestParallelQueue(t *testing.T) {
	assert := assert.New(t)

	wg := sync.WaitGroup{}
	wg.Add(8)
	w := NewParallelQueue(func(_ context.Context, obj interface{}) error {
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
