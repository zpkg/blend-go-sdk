package async

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestAutotriggerAction(t *testing.T) {
	wg := &sync.WaitGroup{}
	wg.Add(1)

	action := func(obj interface{}) {
		wg, ok := obj.(*sync.WaitGroup)
		if ok {
			wg.Done()
		}
	}

	at := NewAutotriggerAction(time.Hour).
		WithHandler(action).
		WithTriggerOnAbort(false)

	at.Start()
	defer at.Stop()

	at.SetValue(wg)
	at.Trigger()
	wg.Wait()
}

func TestAutotriggerActionTicker(t *testing.T) {
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

	at := NewAutotriggerAction(time.Millisecond * 3).
		WithHandler(action).
		WithTriggerOnAbort(false)

	at.Start()
	defer at.Stop()
	wg.Wait()
}
