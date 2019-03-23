package logger

import (
	"net/http"

	"github.com/blend/go-sdk/async"
)

const (

	// DefaultListenerName is a default.
	DefaultListenerName = "default"

	// DefaultRecoverPanics is a default.
	DefaultRecoverPanics = true
)

// New returns a new logger with a given set of enabled flags, without a writer provisioned.
func New(options ...Option) *Logger {
	l := &Logger{
		Latch:         async.NewLatch(),
		RecoverPanics: DefaultRecoverPanics,
		Flags:         NewFlags(),
	}
	return l
}

// Option is a logger option.
type Option func(*Logger) error

// Logger is a handler for various logging events with descendent handlers.
type Logger struct {
	*async.Latch

	RecoverPanics bool
	Writers       []Writer
	Flags         *Flags
	Workers       map[string]map[string]*Worker
	WriteWorker   *Worker
}

// WithEnabled flips the bit flag for a given set of events.
func (l *Logger) WithEnabled(flags ...string) *Logger {
	l.Enable(flags...)
	return l
}

// Enable flips the bit flag for a given set of events.
func (l *Logger) Enable(flags ...string) {
	l.flagsLock.Lock()
	defer l.flagsLock.Unlock()

	if l.flags != nil {
		for _, flag := range flags {
			l.flags.Enable(flag)
		}
	} else {
		l.flags = NewFlagSet(flags...)
	}
}

// WithDisabled flips the bit flag for a given set of events.
func (l *Logger) WithDisabled(flags ...string) *Logger {
	l.Disable(flags...)
	return l
}

// Disable flips the bit flag for a given set of events.
func (l *Logger) Disable(flags ...string) {
	l.flagsLock.Lock()
	defer l.flagsLock.Unlock()
	for _, flag := range flags {
		l.flags.Disable(flag)
	}
}

// IsEnabled asserts if a flag value is set or not.
func (l *Logger) IsEnabled(flag string) (enabled bool) {
	l.flagsLock.Lock()
	if l.flags == nil {
		enabled = false
		l.flagsLock.Unlock()
		return
	}
	enabled = l.flags.IsEnabled(flag)
	l.flagsLock.Unlock()
	return
}

// HasListeners returns if there are registered listener for an event.
func (l *Logger) HasListeners(flag string) bool {
	l.workersLock.Lock()
	defer l.workersLock.Unlock()

	if l.workers == nil {
		return false
	}
	workers, hasWorkers := l.workers[flag]
	if !hasWorkers {
		return false
	}
	return len(workers) > 0
}

// HasListener returns if a specific listener is registerd for a flag.
func (l *Logger) HasListener(flag, listenerName string) bool {
	l.workersLock.Lock()
	defer l.workersLock.Unlock()

	if l.workers == nil {
		return false
	}
	workers, hasWorkers := l.workers[flag]
	if !hasWorkers {
		return false
	}
	_, hasWorker := workers[listenerName]
	return hasWorker
}

// Listen adds a listener for a given flag.
func (l *Logger) Listen(flag, listenerName string, listener Listener) {
	l.workersLock.Lock()
	defer l.workersLock.Unlock()

	if l.workers == nil {
		l.workers = map[string]map[string]*Worker{}
	}

	w := NewWorker(l, listener, l.listenerWorkerQueueDepth)
	if listeners, hasListeners := l.workers[flag]; hasListeners {
		listeners[listenerName] = w
	} else {
		l.workers[flag] = map[string]*Worker{
			listenerName: w,
		}
	}
	w.Start()
}

// RemoveListeners clears *all* listeners for a Flag.
func (l *Logger) RemoveListeners(flag string) {
	l.workersLock.Lock()
	defer l.workersLock.Unlock()

	if l.workers == nil {
		return
	}

	listeners, hasListeners := l.workers[flag]
	if !hasListeners {
		return
	}

	for _, w := range listeners {
		w.Close()
	}

	delete(l.workers, flag)
}

// RemoveListener clears a specific listener for a Flag.
func (l *Logger) RemoveListener(flag, listenerName string) {
	l.workersLock.Lock()
	defer l.workersLock.Unlock()

	if l.workers == nil {
		return
	}

	listeners, hasListeners := l.workers[flag]
	if !hasListeners {
		return
	}

	worker, hasWorker := listeners[listenerName]
	if !hasWorker {
		return
	}

	worker.Close()
	delete(listeners, listenerName)

	if len(listeners) == 0 {
		delete(l.workers, flag)
	}
}

// Trigger fires the listeners for a given event asynchronously.
// The invocations will be queued in a work queue and processed by a fixed worker count.
// There are no order guarantees on when these events will be processed.
// This call will not block on the event listeners.
func (l *Logger) Trigger(e Event) {
	l.trigger(true, e)
}

// SyncTrigger fires the listeners for a given event synchronously.
// The invocations will be triggered immediately, blocking the call.
func (l *Logger) SyncTrigger(e Event) {
	l.trigger(false, e)
}

func (l *Logger) trigger(async bool, e Event) {
	if !async && l.recoverPanics {
		defer func() {
			if r := recover(); r != nil {
				l.Write(Errorf(Fatal, "%+v", r))
			}
		}()
	}

	if async && !l.isStarted() {
		return
	}

	if typed, isTyped := e.(EventEnabled); isTyped && !typed.IsEnabled() {
		return
	}

	flag := e.Flag()
	if l.IsEnabled(flag) {
		if l.heading != "" {
			if typed, isTyped := e.(EventHeadings); isTyped {
				if len(typed.Headings()) > 0 {
					typed.SetHeadings(append([]string{l.heading}, typed.Headings()...)...)
				} else {
					typed.SetHeadings(l.heading)
				}
			}
		}

		var workers map[string]*Worker
		l.workersLock.Lock()
		if l.workers != nil {
			if flagWorkers, hasWorkers := l.workers[flag]; hasWorkers {
				workers = flagWorkers
			}
		}
		l.workersLock.Unlock()

		for _, worker := range workers {
			if async {
				worker.Work <- e
			} else {
				worker.Listener(e)
			}
		}

		// check if the flag is globally hidden from output.
		if l.IsHidden(flag) {
			return
		}

		// check if the event controls if it should be written or not.
		if typed, isTyped := e.(EventWritable); isTyped && !typed.IsWritable() {
			return
		}

		if async && l.writeWorker != nil {
			l.writeWorker.Work <- e
		} else {
			l.Write(e)
		}
	}
}

// --------------------------------------------------------------------------------
// Builtin Flag Handlers (infof, debugf etc.)
// --------------------------------------------------------------------------------

// Infof logs an informational message to the output stream.
func (l *Logger) Infof(format string, args ...interface{}) {
	l.trigger(true, Messagef(Info, format, args...))
}

// Debugf logs a debug message to the output stream.
func (l *Logger) Debugf(format string, args ...interface{}) {
	l.trigger(true, Messagef(Debug, format, args...))
}

// Warningf logs a debug message to the output stream.
func (l *Logger) Warningf(format string, args ...interface{}) {
	l.trigger(false, Errorf(Warning, format, args...))
}

// Warning logs a warning error to std err.
func (l *Logger) Warning(err error) error {
	l.trigger(true, NewErrorEvent(Warning, err))
	return err
}

// WarningWithReq logs a warning error to std err with a request.
func (l *Logger) WarningWithReq(err error, req *http.Request) error {
	l.trigger(true, NewErrorEventWithState(Warning, err, req))
	return err
}

// Errorf writes an event to the log and triggers event listeners.
func (l *Logger) Errorf(format string, args ...interface{}) {
	l.trigger(true, Errorf(Error, format, args...))
}

// Error logs an error to std err.
func (l *Logger) Error(err error) error {
	l.trigger(true, NewErrorEvent(Error, err))
	return err
}

// ErrorWithReq logs an error to std err with a request.
func (l *Logger) ErrorWithReq(err error, req *http.Request) error {
	l.trigger(true, NewErrorEventWithState(Error, err, req))
	return err
}

// Fatalf writes an event to the log and triggers event listeners.
func (l *Logger) Fatalf(format string, args ...interface{}) {
	l.trigger(true, Errorf(Fatal, format, args...))
}

// Fatal logs the result of a panic to std err.
func (l *Logger) Fatal(err error) error {
	l.trigger(true, NewErrorEvent(Fatal, err))
	return err
}

// Write writes an event synchronously to the writer either as a normal even or as an error.
func (l *Logger) Write(e Event) {
	ll := len(l.writers)
	if typed, isTyped := e.(EventError); isTyped && typed.IsError() {
		for index := 0; index < ll; index++ {
			l.writers[index].WriteError(e)
		}
		return
	}
	for index := 0; index < ll; index++ {
		l.writers[index].Write(e)
	}
}

// --------------------------------------------------------------------------------
// finalizers
// --------------------------------------------------------------------------------

// Close releases shared resources for the agent.
func (l *Logger) Close() (err error) {
	l.flagsLock.Lock()
	defer l.flagsLock.Unlock()

	if l.flags != nil {
		l.flags.SetNone()
	}

	l.setStopping()

	l.workersLock.Lock()
	defer l.workersLock.Unlock()

	for _, workers := range l.workers {
		for _, worker := range workers {
			worker.Close()
		}
	}

	for key := range l.workers {
		delete(l.workers, key)
	}
	l.workers = nil

	l.writeWorkerLock.Lock()
	defer l.writeWorkerLock.Unlock()

	l.writeWorker.Close()
	l.writeWorker = nil

	l.setStopped()

	return nil
}

// Drain waits for the agent to finish its queue of events before closing.
func (l *Logger) Drain() error {
	l.workersLock.Lock()
	defer l.workersLock.Unlock()

	l.setStopping()

	for _, workers := range l.workers {
		for _, worker := range workers {
			worker.Drain()
		}
	}

	l.writeWorkerLock.Lock()
	defer l.writeWorkerLock.Unlock()

	if l.writeWorker != nil {
		l.writeWorker.Drain()
	}

	l.setStarted()

	return nil
}
