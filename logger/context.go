package logger

// Context is a logger context.
// It is used to split a logger into functional concerns
// but retain all the underlying machinery of logging.
type Context struct {
	*Logger
	*Context
}


// --------------------------------------------------------------------------------
// Builtin Flag Handlers (infof, debugf etc.)
// --------------------------------------------------------------------------------

// Infof logs an informational message to the output stream.
func (c *Context) Infof(format string, args ...interface{}) {
	l.Trigger(Messagef(Info, format, args...))
}

// Debugf logs a debug message to the output stream.
func (c *Context) Debugf(format string, args ...interface{}) {
	l.Trigger(Messagef(Debug, format, args...))
}

// Warningf logs a debug message to the output stream.
func (c *Context) Warningf(format string, args ...interface{}) {
	l.Trigger(Errorf(Warning, format, args...))
}

// Warning logs a warning error to std err.
func (c *Context) Warning(err error) error {
	l.Trigger(NewErrorEvent(Warning, err))
	return err
}

// WarningWithReq logs a warning error to std err with a request.
func (c *Context) WarningWithReq(err error, req *http.Request) error {
	l.Trigger(NewErrorEventWithState(Warning, err, req))
	return err
}

// Errorf writes an event to the log and triggers event listeners.
func (c *Context) Errorf(format string, args ...interface{}) {
	l.Trigger(Errorf(Error, format, args...))
}

// Error logs an error to std err.
func (c *Context) Error(err error) error {
	l.Trigger(NewErrorEvent(Error, err))
	return err
}

// ErrorWithReq logs an error to std err with a request.
func (c *Context) ErrorWithReq(err error, req *http.Request) error {
	l.Trigger(NewErrorEventWithState(Error, err, req))
	return err
}

// Fatalf writes an event to the log and triggers event listeners.
func (c *Context) Fatalf(format string, args ...interface{}) {
	l.Trigger(Errorf(Fatal, format, args...))
}

// Fatal logs the result of a panic to std err.
func (c *Context) Fatal(err error) error {
	l.Trigger(NewErrorEvent(Fatal, err))
	return err
}

// Write writes an event synchronously to the writer either as a normal even or as an error.
func (c *Context) Write(e Event) error {
	return l.Formatter(l.Output, e)
}