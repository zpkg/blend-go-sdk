package logger

// Listenable is an interface loggers can ascribe to.
type Listenable interface {
	Listen(flag string, label string, listener Listener)
}

// Triggerable is an interface.
type Triggerable interface {
	Trigger(Event)
}

// InfofReceiver is a type that defines Infof.
type InfofReceiver interface {
	Infof(string, ...interface{})
}

// DebugfReceiver is a type that defines Debugf.
type DebugfReceiver interface {
	Debugf(string, ...interface{})
}

// OutputReceiver is an interface
type OutputReceiver interface {
	InfofReceiver
	DebugfReceiver
}

// WarningfReceiver is a type that defines Warningf.
type WarningfReceiver interface {
	Warningf(string, ...interface{})
}

// ErrorfReceiver is a type that defines Errorf.
type ErrorfReceiver interface {
	Errorf(string, ...interface{})
}

// FatalfReceiver is a type that defines Fatalf.
type FatalfReceiver interface {
	Fatalf(string, ...interface{})
}

// ErrorOutputReceiver is an interface
type ErrorOutputReceiver interface {
	WarningfReceiver
	ErrorfReceiver
	FatalfReceiver
}

// WarningReceiver is a type that defines Warning.
type WarningReceiver interface {
	Warning(error) error
}

// ErrorReceiver is a type that defines Error.
type ErrorReceiver interface {
	Error(error) error
}

// FatalReceiver is a type that defines Fatal.
type FatalReceiver interface {
	Fatal(error) error
}

// Errorable is an interface
type Errorable interface {
	WarningReceiver
	ErrorReceiver
	FatalReceiver
}

// Log is a logger that implements the full suite of logging methods.
type Log interface {
	Listenable
	Triggerable
	OutputReceiver
	ErrorOutputReceiver
	Errorable
}
