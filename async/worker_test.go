package async

import (
	"context"
	"sync"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestWorker(t *testing.T) {
	assert := assert.New(t)

	var didWork bool
	wg := sync.WaitGroup{}
	wg.Add(1)
	w := NewWorker(func(_ context.Context, obj interface{}) error {
		defer wg.Done()
		didWork = true
		assert.Equal("hello", obj)
		return nil
	})

	w.Start()
	assert.True(w.Latch().IsRunning())
	w.Enqueue("hello")
	wg.Wait()
	w.Close()
	assert.False(w.Latch().IsRunning())
	assert.True(didWork)
}
