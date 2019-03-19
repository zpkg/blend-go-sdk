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
	w := NewQueue(func(_ context.Context, obj interface{}) error {
		defer wg.Done()
		return nil
	})
	w.Start()
	assert.True(w.IsRunning())

	for x := 0; x < 8; x++ {
		w.Enqueue("hello")
	}

	wg.Wait()
	w.Close()
	assert.False(w.IsRunning())
}
