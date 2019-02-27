package async

import (
	"sync"
	"time"

	"github.com/blend/go-sdk/exception"
)

// NewAutotriggerAction returns a new NewAutotriggerAction
func NewAutotriggerAction(interval time.Duration) *AutotriggerAction {
	return &AutotriggerAction{
		value:          nil,
		valueLock:      &sync.Mutex{},
		handler:        nil,
		interval:       interval,
		latch:          NewLatch(),
		triggerOnAbort: true,
	}
}

// AutotriggerAction is an action that is triggered automatically on some set interval.
// It also exposes a function to trigger the action synchronously
type AutotriggerAction struct {
	value          interface{}
	valueLock      *sync.Mutex
	handler        func(interface{})
	interval       time.Duration
	latch          *Latch
	triggerOnAbort bool
}

// WithHandler sets the trigger handler
func (at *AutotriggerAction) WithHandler(handler func(interface{})) *AutotriggerAction {
	at.handler = handler
	return at
}

// WithTriggerOnAbort determines whether the trigger should be invoked when AutotriggerAction is stopped
func (at *AutotriggerAction) WithTriggerOnAbort(triggerOnAbort bool) *AutotriggerAction {
	at.triggerOnAbort = triggerOnAbort
	return at
}

// NotifyStarted returns the started signal.
func (at *AutotriggerAction) NotifyStarted() <-chan struct{} {
	return at.latch.NotifyStarted()
}

// NotifyStopped returns the started stopped.
func (at *AutotriggerAction) NotifyStopped() <-chan struct{} {
	return at.latch.NotifyStopped()
}

// Start starts the trigger.
func (at *AutotriggerAction) Start() error {
	if !at.latch.CanStart() {
		return exception.New(ErrCannotStart)
	}
	at.latch.Starting()
	go func() {
		at.latch.Started()
		at.runLoop()
	}()
	<-at.latch.NotifyStarted()
	return nil
}

// Stop stops the trigger
func (at *AutotriggerAction) Stop() error {
	if !at.latch.CanStop() {
		return exception.New(ErrCannotStop)
	}
	at.latch.Stopping()
	<-at.latch.NotifyStopped()
	return nil
}

func (at *AutotriggerAction) runLoop() {
	ticker := time.Tick(at.interval)
	for {
		select {
		case <-ticker:
			at.TriggerAsync()
		case <-at.latch.NotifyStopping():
			if at.triggerOnAbort {
				at.Trigger()
			}
			at.latch.Stopped()
			return
		}
	}
}

// Set sets the value
func (at *AutotriggerAction) SetValue(value interface{}) {
	at.valueLock.Lock()
	at.value = value
	defer at.valueLock.Unlock()
}

// Trigger invokes the handler, if one is set, with the value
// This call is synchronous, in that it will call the trigger handler on the same goroutine.
func (at *AutotriggerAction) Trigger() {
	at.valueLock.Lock()
	defer at.valueLock.Unlock()
	at.triggerUnsafe()
}

// TriggerAsync calls the handler, if one is set, with the value
// This call is asynchronous, in that it will call the trigger handler on its own goroutine.
func (at *AutotriggerAction) TriggerAsync() {
	at.valueLock.Lock()
	defer at.valueLock.Unlock()
	at.triggerUnsafeAsync()
}

// triggerUnsafe calls the handler, if one is set, without acquiring any locks
func (at *AutotriggerAction) triggerUnsafe() {
	if at.handler != nil {
		at.handler(at.value)
	}
}

// triggerUnsafeAsync calls the handler, if one is set, without acquiring any locks
func (at *AutotriggerAction) triggerUnsafeAsync() {
	if at.handler != nil {
		go at.handler(at.value)
	}
}
