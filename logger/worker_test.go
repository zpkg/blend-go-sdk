package logger

import (
	"bytes"
	"sync"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestWorker(t *testing.T) {
	assert := assert.New(t)

	wg := sync.WaitGroup{}
	wg.Add(1)
	var didFire bool
	w := NewWorker(func(e Event) {
		defer wg.Done()
		didFire = true

		typed, isTyped := e.(*MessageEvent)
		assert.True(isTyped)
		assert.Equal("test", typed.Message)
	})

	w.Start()
	defer w.Stop()

	w.Work <- Messagef(Info, "test")
	wg.Wait()

	assert.True(didFire)
}

func TestWorkerPanics(t *testing.T) {
	assert := assert.New(t)

	buffer := bytes.NewBuffer(nil)

	log := New(OptAll(), OptOutput(buffer), OptFormatter(&TextFormatter{}))
	defer log.Close()

	wg := sync.WaitGroup{}
	wg.Add(1)
	var didFire bool
	w := NewWorker(func(e Event) {
		defer wg.Done()
		didFire = true
		panic("only a test")
	})
	w.Start()

	w.Work <- Messagef(Info, "test")
	wg.Wait()

	assert.True(didFire)
	w.Stop()
	assert.NotEmpty(buffer.String())
}

func TestWorkerDrain(t *testing.T) {
	assert := assert.New(t)

	wg := sync.WaitGroup{}
	wg.Add(4)
	var didFire bool
	w := NewWorker(func(e Event) {
		defer wg.Done()
		didFire = true
	})

	w.Work <- Messagef(Info, "test1")
	w.Work <- Messagef(Info, "test2")
	w.Work <- Messagef(Info, "test3")
	w.Work <- Messagef(Info, "test4")

	go func() {
		w.Drain()
	}()
	wg.Wait()

	assert.True(didFire)
}
