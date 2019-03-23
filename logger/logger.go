package logger

import (
	"io"
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

// Logger is a handler for various logging events with descendent handlers.
type Logger struct {
	*async.Latch
	*Flags

	RecoverPanics bool

	Output    io.Writer
	Errors    chan error
	Listeners map[string]map[string]*Worker
}

// HasListeners returns if there are registered listener for an event.
func (l *Logger) HasListeners(flag string) bool {
	l.Lock()
	defer l.Unlock()

	if l.Listeners == nil {
		return false
	}
	listeners, ok := l.Listeners[flag]
	if !ok {
		return false
	}
	return len(listeners) > 0
}

// HasListener returns if a specific listener is registerd for a flag.
func (l *Logger) HasListener(flag, listenerName string) bool {
	l.Lock()
	defer l.Unlock()

	if l.Listeners == nil {
		return false
	}
	workers, ok := l.Listeners[flag]
	if !ok {
		return false
	}
	_, ok = workers[listenerName]
	return ok
}

// Listen adds a listener for a given flag.
func (l *Logger) Listen(flag, listenerName string, listener Listener) {
	l.Lock()
	defer l.Unlock()

	if l.Listeners == nil {
		l.Listeners = make(map[string]map[string]*Worker)
	}

	w := NewWorker(listener)
	if listeners, ok := l.Listeners[flag]; ok {
		listeners[listenerName] = w
	} else {
		l.Listeners[flag] = map[string]*Worker{
			listenerName: w,
		}
	}
	go w.Start()
	<-w.NotifyStarted()
}

// RemoveListeners clears *all* listeners for a Flag.
func (l *Logger) RemoveListeners(flag string) {
	l.Lock()
	defer l.Unlock()

	if l.Listeners == nil {
		return
	}

	listeners, ok := l.Listeners[flag]
	if !ok {
		return
	}

	for _, l := range listeners {
		l.Stop()
	}

	delete(l.Listeners, flag)
}

// RemoveListener clears a specific listener for a Flag.
func (l *Logger) RemoveListener(flag, listenerName string) {
	l.Lock()
	defer l.Unlock()

	if l.Listeners == nil {
		return
	}

	listeners, ok := l.Listeners[flag]
	if !ok {
		return
	}

	worker, ok := listeners[listenerName]
	if !ok {
		return
	}

	worker.Stop()
	<-worker.NotifyStopped()

	delete(listeners, listenerName)
	if len(listeners) == 0 {
		delete(l.Listeners, flag)
	}
}

// Trigger fires the listeners for a given event asynchronously.
// The invocations will be queued in a work queue per listener.
// There are no order guarantees on when these events will be processed across listeners.
// This call will not block on the event listeners, but will block on writing the event to the formatted output.
func (l *Logger) Trigger(e Event) {
	flag := e.Flag()
	if !l.IsEnabled(flag) {
		return
	}

	if typed, isTyped := e.(EnabledProvider); isTyped && !typed.IsEnabled() {
		return
	}

	var listeners map[string]*Worker
	l.Lock()
	if l.Listeners != nil {
		if flagListeners, ok := l.Listeners[flag]; ok {
			listeners = flagListeners
		}
	}
	l.Unlock()

	for _, listener := range listeners {
		listener.Work <- e
	}

	// check if the event controls if it should be written or not.
	if typed, isTyped := e.(WritableProvider); isTyped && !typed.IsWritable() {
		return
	}

	if err := l.Write(e); err != nil && l.Errors != nil {
		l.Errors <- err
	}
}

// --------------------------------------------------------------------------------
// Builtin Flag Handlers (infof, debugf etc.)
// --------------------------------------------------------------------------------

// Infof logs an informational message to the output stream.
func (l *Logger) Infof(format string, args ...interface{}) {
	l.Trigger(Messagef(Info, format, args...))
}

// Debugf logs a debug message to the output stream.
func (l *Logger) Debugf(format string, args ...interface{}) {
	l.Trigger(Messagef(Debug, format, args...))
}

// Warningf logs a debug message to the output stream.
func (l *Logger) Warningf(format string, args ...interface{}) {
	l.Trigger(Errorf(Warning, format, args...))
}

// Warning logs a warning error to std err.
func (l *Logger) Warning(err error) error {
	l.Trigger(NewErrorEvent(Warning, err))
	return err
}

// WarningWithReq logs a warning error to std err with a request.
func (l *Logger) WarningWithReq(err error, req *http.Request) error {
	l.Trigger(NewErrorEventWithState(Warning, err, req))
	return err
}

// Errorf writes an event to the log and triggers event listeners.
func (l *Logger) Errorf(format string, args ...interface{}) {
	l.Trigger(Errorf(Error, format, args...))
}

// Error logs an error to std err.
func (l *Logger) Error(err error) error {
	l.Trigger(NewErrorEvent(Error, err))
	return err
}

// ErrorWithReq logs an error to std err with a request.
func (l *Logger) ErrorWithReq(err error, req *http.Request) error {
	l.Trigger(NewErrorEventWithState(Error, err, req))
	return err
}

// Fatalf writes an event to the log and triggers event listeners.
func (l *Logger) Fatalf(format string, args ...interface{}) {
	l.Trigger(Errorf(Fatal, format, args...))
}

// Fatal logs the result of a panic to std err.
func (l *Logger) Fatal(err error) error {
	l.Trigger(NewErrorEvent(Fatal, err))
	return err
}

// Write writes an event synchronously to the writer either as a normal even or as an error.
func (l *Logger) Write(e Event) error {
	return l.Formatter(l.Output, e)
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
