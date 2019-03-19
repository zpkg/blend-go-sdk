package async

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestAutoActionTrigger(t *testing.T) {
	wg := &sync.WaitGroup{}
	wg.Add(1)

	action := func(_ context.Context) error {
		wg.Done()
		return nil
	}

	at := NewAutoAction(action, time.Hour, 10)
	go at.Start()
	<-at.NotifyStarted()

	defer at.Stop()
	at.Trigger(context.Background())
	wg.Wait()
}

func TestAutoActionTick(t *testing.T) {
	wg := sync.WaitGroup{}

	// keep track of the wait group state
	ticksRemaining := int32(3)
	wg.Add(int(ticksRemaining))

	action := func(_ context.Context) error {
		ticks := atomic.LoadInt32(&ticksRemaining)
		if ticks > 0 {
			atomic.AddInt32(&ticksRemaining, -1)
			wg.Done()
		}
		return nil
	}

	at := NewAutoAction(action, time.Millisecond, 10)
	go at.Start()
	<-at.NotifyStarted()
	defer at.Stop()
	wg.Wait()
}

func TestAutoActionCount(t *testing.T) {
	wg := sync.WaitGroup{}

	wg.Add(1)

	action := func(_ context.Context) error {
		wg.Done()
		return nil
	}

	maxCounter := 10

	at := NewAutoAction(action, time.Hour, maxCounter)

	for i := 0; i < maxCounter; i++ {
		at.Increment(context.Background())
	}

	go at.Start()
	<-at.NotifyStarted()
	defer at.Stop()
	wg.Wait()
}
