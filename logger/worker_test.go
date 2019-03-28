package logger

import (
	"context"
	"sync"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestWorker(t *testing.T) {
	assert := assert.New(t)

	wg := sync.WaitGroup{}
	wg.Add(1)
	var didFire bool
	w := NewWorker(func(_ context.Context, e Event) {
		defer wg.Done()
		didFire = true

		typed, isTyped := e.(*MessageEvent)
		assert.True(isTyped)
		assert.Equal("test", typed.Message)
	})

	go w.Start()
	<-w.NotifyStarted()
	defer w.Stop()

	w.Work <- EventWithContext{context.Background(), Messagef(Info, "test")}
	wg.Wait()

	assert.True(didFire)
}

func TestWorkerDrain(t *testing.T) {
	assert := assert.New(t)

	wg := sync.WaitGroup{}
	wg.Add(4)
	var didFire bool
	w := NewWorker(func(ctx context.Context, e Event) {
		defer wg.Done()
		didFire = true
	})

	w.Work <- EventWithContext{context.Background(), Messagef(Info, "test1")}
	w.Work <- EventWithContext{context.Background(), Messagef(Info, "test2")}
	w.Work <- EventWithContext{context.Background(), Messagef(Info, "test3")}
	w.Work <- EventWithContext{context.Background(), Messagef(Info, "test4")}

	go func() {
		w.Drain()
	}()
	wg.Wait()

	assert.True(didFire)
}
