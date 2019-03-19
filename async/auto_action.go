package async

import (
	"sync"
	"sync/atomic"
	"time"

	"github.com/blend/go-sdk/exception"
)

// NewAutoAction returns a new NewAutoAction
func NewAutoAction(interval time.Duration, maxCount int32) *AutoAction {
	return &AutoAction{
		maxCount:       maxCount,
		interval:       interval,
		latch:          NewLatch(),
		triggerOnAbort: true,
	}
}

// AutoAction is an action that is triggered automatically on some set interval.
// It also exposes a function to trigger the action synchronously
type AutoAction struct {
	sync.Mutex
	counter        int32
	maxCount       int32
	action         func()
	interval       time.Duration
	latch          *Latch
	triggerOnAbort bool
}

// MaxCount returns the number of increments between action triggers
func (a *AutoAction) MaxCount() int32 {
	return a.maxCount
}

// WithAction sets the trigger action
// This should be called before Start(), and, if it's called after start, the AutoAction's lock should be acquired first
func (a *AutoAction) WithAction(action func()) *AutoAction {
	a.action = action
	return a
}

// ShouldTriggerOnAbort refers to whether the action should be triggered when NewAutoAction is stopped
func (a *AutoAction) ShouldTriggerOnAbort() bool {
	return a.triggerOnAbort
}

// WithTriggerOnAbort sets whether the action should be triggered when NewAutoAction is stopped
func (a *AutoAction) WithTriggerOnAbort(triggerOnAbort bool) *AutoAction {
	a.triggerOnAbort = triggerOnAbort
	return a
}

// NotifyStarted returns the started signal.
func (a *AutoAction) NotifyStarted() <-chan struct{} {
	return a.latch.NotifyStarted()
}

// NotifyStopped returns the started stopped.
func (a *AutoAction) NotifyStopped() <-chan struct{} {
	return a.latch.NotifyStopped()
}

// Start starts the trigger.
func (a *AutoAction) Start() error {
	if !a.latch.CanStart() {
		return exception.New(ErrCannotStart)
	}
	a.latch.Starting()
	go func() {
		a.latch.Started()
		a.runLoop()
	}()
	<-a.latch.NotifyStarted()
	return nil
}

// Stop stops the trigger
func (a *AutoAction) Stop() error {
	if !a.latch.CanStop() {
		return exception.New(ErrCannotStop)
	}
	a.latch.Stopping()
	<-a.latch.NotifyStopped()
	return nil
}

func (a *AutoAction) runLoop() {
	ticker := time.Tick(a.interval)
	for {
		select {
		case <-ticker:
			a.Trigger()
		case <-a.latch.NotifyStopping():
			if a.triggerOnAbort {
				a.Trigger()
			}
			a.latch.Stopped()
			return
		}
	}
}

// Increment updates the count
func (a *AutoAction) Increment() {
	if atomic.CompareAndSwapInt32(&a.counter, a.maxCount-1, 0) {
		a.Trigger()
		return
	}
	atomic.AddInt32(&a.counter, 1)
}

// Trigger invokes the action, if one is set, with the value
// This call is synchronous, in that it will call the trigger action on the same goroutine.
func (a *AutoAction) Trigger() {
	a.Lock()
	defer a.Unlock()
	if a.action != nil {
		a.action()
	}
}
