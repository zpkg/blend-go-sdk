package async

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestAutoActionTrigger(t *testing.T) {
	wg := &sync.WaitGroup{}
	wg.Add(1)

	action := func() {
		wg.Done()
	}

	a := NewAutoAction(time.Hour, 10).
		WithAction(action).
		WithTriggerOnAbort(false)

	a.Start()
	defer a.Stop()

	a.Trigger()
	wg.Wait()
}

func TestAutoActionTick(t *testing.T) {
	wg := sync.WaitGroup{}

	// keep track of the wait group state
	ticksRemaining := int32(3)
	wg.Add(int(ticksRemaining))

	action := func() {
		ticks := atomic.LoadInt32(&ticksRemaining)
		if ticks > 0 {
			atomic.AddInt32(&ticksRemaining, -1)
			wg.Done()
		}
	}

	at := NewAutoAction(time.Millisecond*3, 10).
		WithAction(action).
		WithTriggerOnAbort(false)

	at.Start()
	defer at.Stop()

	wg.Wait()
}

func TestAutoActionCount(t *testing.T) {
	wg := sync.WaitGroup{}

	wg.Add(1)

	action := func() {
		wg.Done()
	}

	maxCounter := 10

	at := NewAutoAction(time.Hour, int32(maxCounter)).
		WithAction(action).
		WithTriggerOnAbort(false)

	for i := 0; i < maxCounter; i++ {
		at.Increment()
	}

	at.Start()
	defer at.Stop()
	wg.Wait()
}
