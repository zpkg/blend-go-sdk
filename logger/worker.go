package logger

import (
	"github.com/blend/go-sdk/async"
	"github.com/blend/go-sdk/exception"
)

// NewWorker returns a new worker.
func NewWorker(listener Listener) *Worker {
	return &Worker{
		Latch:    async.NewLatch(),
		Listener: listener,
		Work:     make(chan Event, DefaultWorkerQueueDepth),
	}
}

// Worker is an agent that processes a listener.
type Worker struct {
	*async.Latch
	Errors   chan error
	Listener Listener
	Work     chan Event
}

// Start starts the worker.
func (w *Worker) Start() error {
	if !w.CanStart() {
		return exception.New(async.ErrCannotStart)
	}
	w.Starting()
	w.Dispatch()
	return nil
}

// Dispatch is the main listen loop
func (w *Worker) Dispatch() {
	w.Started()
	var e Event
	var err error
	for {
		select {
		case e = <-w.Work:
			if err = w.Process(e); err != nil && w.Errors != nil {
				w.Errors <- err
			}
		case <-w.NotifyPausing():
			w.Paused()
			<-w.NotifyResuming()
			w.Started()
		case <-w.NotifyStopping():
			w.Stopped()
			return
		}
	}
}

// Process calls the listener for an event.
func (w *Worker) Process(e Event) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = exception.New(r)
			return
		}
	}()
	w.Listener(e)
	return
}

// Drain stops the worker and synchronously processes any remaining work.
// It then restarts the worker.
func (w *Worker) Drain() {
	w.Pausing()
	<-w.NotifyPaused()

	var work Event
	var err error
	workLeft := len(w.Work)
	for index := 0; index < workLeft; index++ {
		work = <-w.Work
		if err = w.Process(work); err != nil && w.Errors != nil {
			w.Errors <- err
		}
	}

	w.Resuming()
	<-w.NotifyStarted()
}

// Stop stops the worker.
func (w *Worker) Stop() error {
	if !w.CanStop() {
		return exception.New(async.ErrCannotStop)
	}
	w.Stopping()
	<-w.NotifyStopped()

	var work Event
	var err error

	workLeft := len(w.Work)
	for index := 0; index < workLeft; index++ {
		work = <-w.Work
		if err = w.Process(work); err != nil && w.Errors != nil {
			w.Errors <- err
		}
	}
	return nil
}
