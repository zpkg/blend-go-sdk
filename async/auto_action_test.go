package async

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestAutoAction(t *testing.T) {
	wg := &sync.WaitGroup{}
	wg.Add(1)

	action := func(obj interface{}) {
		wg, ok := obj.(*sync.WaitGroup)
		if ok {
			wg.Done()
		}
	}

	a := NewAutoAction(time.Hour, 10).
		WithHandler(action).
		WithTriggerOnAbort(false)

	a.Start()
	defer a.Stop()

	a.Update(wg)
	a.Trigger()
	wg.Wait()
}

func TestAutoActionTick(t *testing.T) {
	wg := sync.WaitGroup{}

	// keep track of the wait group state
	ticksRemaining := int32(3)
	wg.Add(int(ticksRemaining))

	action := func(obj interface{}) {
		ticks := atomic.LoadInt32(&ticksRemaining)
		if ticks > 0 {
			atomic.AddInt32(&ticksRemaining, -1)
			wg.Done()
		}
	}

	at := NewAutoAction(time.Millisecond*3, 10).
		WithHandler(action).
		WithTriggerOnAbort(false)

	at.Start()
	defer at.Stop()
	wg.Wait()
}

func TestAutoActionCount(t *testing.T) {
	wg := sync.WaitGroup{}

	wg.Add(1)

	action := func(obj interface{}) {
		wg.Done()
	}

	at := NewAutoAction(time.Hour, 1).
		WithHandler(action).
		WithTriggerOnAbort(false)
	at.Update(nil)

	at.Start()
	defer at.Stop()
	wg.Wait()
}
