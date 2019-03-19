package async

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/graceful"
)

// Assert AutoflushBuffer is graceful.
var (
	_ graceful.Graceful = (*AutoflushBuffer)(nil)
)

func TestAutoflushBuffer(t *testing.T) {
	assert := assert.New(t)

	wg := sync.WaitGroup{}
	wg.Add(2)

	buffer := NewAutoflushBuffer(func(_ context.Context, objects []interface{}) error {
		defer wg.Done()
		assert.Len(objects, 10)
		return nil
	}, OptAutoflushBufferMaxLen(10), OptAutoflushBufferInterval(time.Hour))

	go buffer.Start()
	<-buffer.NotifyStarted()

	defer buffer.Stop()

	for x := 0; x < 20; x++ {
		buffer.Add(fmt.Sprintf("foo%d", x))
	}

	wg.Wait()
}

func TestAutoflushBufferTicker(t *testing.T) {
	assert := assert.New(t)
	assert.StartTimeout(500 * time.Millisecond)
	defer assert.EndTimeout()

	wg := sync.WaitGroup{}
	wg.Add(20)
	buffer := NewAutoflushBuffer(func(_ context.Context, objects []interface{}) error {
		for range objects {
			wg.Done()
		}
		return nil
	}, OptAutoflushBufferMaxLen(100), OptAutoflushBufferInterval(time.Millisecond))

	go buffer.Start()
	<-buffer.NotifyStarted()

	defer buffer.Stop()

	for x := 0; x < 20; x++ {
		buffer.Add(fmt.Sprintf("foo%d", x))
	}
	wg.Wait()
}

func BenchmarkAutoflushBuffer(b *testing.B) {
	buffer := NewAutoflushBuffer(func(_ context.Context, objects []interface{}) error {
		if len(objects) > 128 {
			b.Fail()
		}
		return nil
	}, OptAutoflushBufferMaxLen(128), OptAutoflushBufferInterval(500*time.Millisecond))

	go buffer.Start()
	<-buffer.NotifyStarted()

	defer buffer.Stop()

	for x := 0; x < b.N; x++ {
		for y := 0; y < 1000; y++ {
			buffer.Add(fmt.Sprintf("asdf%d%d", x, y))
		}
	}
}
