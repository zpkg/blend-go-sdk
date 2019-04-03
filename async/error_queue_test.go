package async

import (
	"context"
	"fmt"
	"sync"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestErrorEqueue(t *testing.T) {
	assert := assert.New(t)

	wg := sync.WaitGroup{}
	wg.Add(5)

	eq := NewErrorQueue(func(_ context.Context, err error) {
		defer wg.Done()
	})
	// you must fire up the worker before you use the work channel.
	// as it gets overwritten on start.
	go eq.Start()
	<-eq.NotifyStarted()

	q := NewQueue(func(_ context.Context, obj interface{}) error {
		return fmt.Errorf("only a test")
	}, OptQueueErrors(eq.Work))

	go q.Start()
	<-q.NotifyStarted()

	for x := 0; x < 5; x++ {
		q.Enqueue("hello")
	}

	wg.Wait()
	eq.Close()
	q.Close()

	assert.False(eq.IsStarted())
	assert.False(q.IsStarted())
}
