package worker

import "time"

// NewInterval returns a new worker that runs an action on an interval.
func NewInterval(action func() error, interval time.Duration) *Interval {
	return &Interval{
		interval: interval,
		action:   action,
		latch:    &Latch{},
	}
}

// Interval is a managed goroutine that does things.
type Interval struct {
	interval time.Duration
	action   func() error
	latch    *Latch
	errors   chan error
}

// Interval returns the interval for the ticker.
func (i Interval) Interval() time.Duration {
	return i.interval
}

// Latch returns the inteval worker latch.
func (i *Interval) Latch() *Latch {
	return i.latch
}

// WithAction sets the interval action.
func (i *Interval) WithAction(action func() error) *Interval {
	i.action = action
	return i
}

// Action returns the interval action.
func (i *Interval) Action() func() error {
	return i.action
}

// WithErrors returns the error channel.
func (i *Interval) WithErrors(errors chan error) *Interval {
	i.errors = errors
	return i
}

// Errors returns a channel to read action errors from.
func (i *Interval) Errors() chan error {
	return i.errors
}

// Start starts the worker.
func (i *Interval) Start() {
	i.latch.Starting()
	go func() {
		i.latch.Started()
		tick := time.Tick(i.interval)
		var err error
		for {
			select {
			case <-tick:
				err = i.action()
				if err != nil && i.errors != nil {
					i.errors <- err
				}
			case <-i.latch.NotifyStop():
				i.latch.Stopped()
				return
			}
		}
	}()
	<-i.latch.NotifyStarted()
}

// Stop stops the worker.
func (i *Interval) Stop() {
	i.latch.Stop()
	<-i.latch.NotifyStopped()
}
