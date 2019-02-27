package async

import (
	"sync"
	"time"

	"github.com/blend/go-sdk/exception"
)

// NewAutoAction returns a new NewAutoAction
func NewAutoAction(interval time.Duration) *AutoAction {
	return &AutoAction{
		value:          nil,
		valueLock:      &sync.Mutex{},
		handler:        nil,
		interval:       interval,
		latch:          NewLatch(),
		triggerOnAbort: true,
	}
}

// NewAutoAction is an action that is triggered automatically on some set interval.
// It also exposes a function to trigger the action synchronously
type AutoAction struct {
	value          interface{}
	valueLock      *sync.Mutex
	handler        func(interface{})
	interval       time.Duration
	latch          *Latch
	triggerOnAbort bool
}

// WithHandler sets the trigger handler
func (a *AutoAction) WithHandler(handler func(interface{})) *AutoAction {
	a.handler = handler
	return a
}

// WithTriggerOnAbort determines whether the action should be triggered when NewAutoAction is stopped
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
			a.TriggerAsync()
		case <-a.latch.NotifyStopping():
			if a.triggerOnAbort {
				a.Trigger()
			}
			a.latch.Stopped()
			return
		}
	}
}

// Set sets the value
func (a *AutoAction) SetValue(value interface{}) {
	a.valueLock.Lock()
	a.value = value
	defer a.valueLock.Unlock()
}

// Trigger invokes the handler, if one is set, with the value
// This call is synchronous, in that it will call the trigger handler on the same goroutine.
func (a *AutoAction) Trigger() {
	a.valueLock.Lock()
	defer a.valueLock.Unlock()
	a.triggerUnsafe()
}

// TriggerAsync calls the handler, if one is set, with the value
// This call is asynchronous, in that it will call the trigger handler on its own goroutine.
func (a *AutoAction) TriggerAsync() {
	a.valueLock.Lock()
	defer a.valueLock.Unlock()
	a.triggerUnsafeAsync()
}

// triggerUnsafe calls the handler, if one is set, without acquiring any locks
func (a *AutoAction) triggerUnsafe() {
	if a.handler != nil {
		a.handler(a.value)
	}
}

// triggerUnsafeAsync calls the handler, if one is set, without acquiring any locks
func (a *AutoAction) triggerUnsafeAsync() {
	if a.handler != nil {
		go a.handler(a.value)
	}
}
