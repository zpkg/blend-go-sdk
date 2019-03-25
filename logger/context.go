package logger

import (
	"context"
	"net/http"
)

// NewContext returns a new context.
func NewContext(log WriteTriggerable, path ...string) *Context {
	return &Context{
		WriteTriggerable: log,
		SubContextPath:   path,
	}
}

// Context is a logger context.
// It is used to split a logger into functional concerns
// but retain all the underlying machinery of logging.
type Context struct {
	WriteTriggerable
	Context        context.Context
	SubContextPath []string
}

// SubContext returns a new sub context.
func (sc *Context) SubContext(name string) *Context {
	return NewContext(sc.WriteTriggerable, append(sc.SubContextPath, name)...)
}

// Background returns the background context.
func (sc *Context) Background() context.Context {
	return WithSubContextPath(sc.Context, sc.SubContextPath)
}

// --------------------------------------------------------------------------------
// Builtin Flag Handlers (infof, debugf etc.)
// --------------------------------------------------------------------------------

// Infof logs an informational message to the output stream.
func (sc *Context) Infof(format string, args ...interface{}) {
	sc.handle(Messagef(Info, format, args...))
}

// Debugf logs a debug message to the output stream.
func (sc *Context) Debugf(format string, args ...interface{}) {
	sc.handle(Messagef(Debug, format, args...))
}

// Warningf logs a debug message to the output stream.
func (sc *Context) Warningf(format string, args ...interface{}) {
	sc.handle(Errorf(Warning, format, args...))
}

// Errorf writes an event to the log and triggers event listeners.
func (sc *Context) Errorf(format string, args ...interface{}) {
	sc.handle(Errorf(Error, format, args...))
}

// Fatalf writes an event to the log and triggers event listeners.
func (sc *Context) Fatalf(format string, args ...interface{}) {
	sc.handle(Errorf(Fatal, format, args...))
}

// Warning logs a warning error to std err.
func (sc *Context) Warning(err error) error {
	sc.handle(NewErrorEvent(Warning, err))
	return err
}

// WarningWithReq logs a warning error to std err with a request.
func (sc *Context) WarningWithReq(err error, req *http.Request) error {
	ee := NewErrorEvent(Warning, err)
	ee.State = req
	sc.handle(ee)
	return err
}

// Error logs an error to std err.
func (sc *Context) Error(err error) error {
	sc.handle(NewErrorEvent(Error, err))
	return err
}

// ErrorWithReq logs an error to std err with a request.
func (sc *Context) ErrorWithReq(err error, req *http.Request) error {
	ee := NewErrorEvent(Error, err)
	ee.State = req
	sc.handle(ee)
	return err
}

// Fatal logs the result of a panic to std err.
func (sc *Context) Fatal(err error) error {
	sc.handle(NewErrorEvent(Fatal, err))
	return err
}

func (sc *Context) handle(e Event) {
	sc.Trigger(sc.Background(), e)
	sc.Write(sc.Background(), e)
}
