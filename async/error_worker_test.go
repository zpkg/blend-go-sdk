package async

import (
	"context"
	"fmt"
	"sync"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestErrorWorker(t *testing.T) {
	assert := assert.New(t)

	var didWork bool
	wg := sync.WaitGroup{}
	wg.Add(1)
	ew := NewErrorWorker(func(_ context.Context, obj error) error {
		defer wg.Done()
		didWork = true
		assert.Equal("hello", obj.Error())
		return nil
	})

	ew.Start()
	assert.True(ew.Latch().IsRunning())
	ew.Enqueue(fmt.Errorf("hello"))
	wg.Wait()
	ew.Close()
	assert.False(ew.Latch().IsRunning())
	assert.True(didWork)
}
