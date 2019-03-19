package async

import (
	"context"
	"sync/atomic"
	"time"

	"github.com/blend/go-sdk/exception"
)

// NewAutoAction returns a new singleton that will trigger an action after a given amount of time has passed (interval)
// or after a given number of increments has happened (maxCount).
func NewAutoAction(action ContextAction, interval time.Duration, maxCount int, options ...AutoActionOption) *AutoAction {
	aa := AutoAction{
		Action:   action,
		MaxCount: int32(maxCount),
		Interval: interval,
		Context:  context.Background(),
	}
	for _, option := range options {
		option(&aa)
	}
	return &aa
}

// AutoActionOption is an option for an auto-action.
type AutoActionOption func(*AutoAction)

// OptAutoActionMaxCount sets the auto-action max count.
func OptAutoActionMaxCount(maxCount int32) AutoActionOption {
	return func(aa *AutoAction) {
		aa.MaxCount = maxCount
	}
}

// OptAutoActionInterval sets the auto-action interval.
func OptAutoActionInterval(d time.Duration) AutoActionOption {
	return func(aa *AutoAction) {
		aa.Interval = d
	}
}

// OptAutoActionErrors sets the auto-action error channel.
func OptAutoActionErrors(errors chan error) AutoActionOption {
	return func(aa *AutoAction) {
		aa.Errors = errors
	}
}

// OptAutoActionTriggerOnStop sets if the auto-action should call the action on shutdown.
func OptAutoActionTriggerOnStop(errors chan error) AutoActionOption {
	return func(aa *AutoAction) {
		aa.Errors = errors
	}
}

// AutoAction is an action that is triggered automatically on some set interval.
// It also exposes a function to trigger the action synchronously
type AutoAction struct {
	Latch

	Action         ContextAction
	Context        context.Context
	Errors         chan error
	Interval       time.Duration
	MaxCount       int32
	TriggerOnAbort bool

	Counter int32
}

// Background returns a background context.
func (a *AutoAction) Background() context.Context {
	if a.Context != nil {
		return a.Context
	}
	return context.Background()
}

/*
Start starts the singleton.

This call blocks. To call it asynchronously:

	go a.Start()
	<-a.NotifyStarted()

This will start the singleton and wait for it to enter the running state.
*/
func (a *AutoAction) Start() error {
	if !a.CanStart() {
		return exception.New(ErrCannotStart)
	}
	a.Starting()
	a.Dispatch()
	return nil
}

// Stop stops the auto-action singleton.
func (a *AutoAction) Stop() error {
	if !a.CanStop() {
		return exception.New(ErrCannotStop)
	}
	a.Stopping()
	<-a.NotifyStopped()
	return nil
}

// Dispatch is the main run loop.
func (a *AutoAction) Dispatch() {
	a.Started()
	ticker := time.Tick(a.Interval)
	for {
		select {
		case <-ticker:
			a.Trigger(a.Background())
		case <-a.NotifyStopping():
			if a.TriggerOnAbort {
				a.Trigger(a.Background())
			}
			a.Stopped()
			return
		}
	}
}

// Increment updates the count
func (a *AutoAction) Increment(ctx context.Context) {
	if atomic.CompareAndSwapInt32(&a.Counter, a.MaxCount-1, 0) {
		a.Trigger(ctx)
		return
	}
	atomic.AddInt32(&a.Counter, 1)
}

// Trigger invokes the action if one is set, it will acquire the lock and hold it for the duration of the call to the action.
func (a *AutoAction) Trigger(ctx context.Context) {
	a.Lock()
	defer a.Unlock()
	defer func() {
		if r := recover(); r != nil {
			if a.Errors != nil {
				a.Errors <- exception.New(r)
			}
		}
	}()

	if err := a.Action(ctx); err != nil && a.Errors != nil {
		a.Errors <- err
	}
}
