package logger

// Listenable is an interface.
type Listenable interface {
	Listen(Flag, string, Listener)
}

// Triggerable is an interface.
type Triggerable interface {
	Trigger(Event)
}

// SyncTriggerable is an interface.
type SyncTriggerable interface {
	SyncTrigger(Event)
}

// OutputReceiver is an interface
type OutputReceiver interface {
	Infof(string, ...Any)
	Sillyf(string, ...Any)
	Debugf(string, ...Any)
}

// SyncOutputReceiver is an interface
type SyncOutputReceiver interface {
	SyncInfof(string, ...Any)
	SyncSillyf(string, ...Any)
	SyncDebugf(string, ...Any)
}

// ErrorOutputReceiver is an interface
type ErrorOutputReceiver interface {
	Warningf(string, ...Any)
	Errorf(string, ...Any)
	Fatalf(string, ...Any)
}

// SyncErrorOutputReceiver is an interface
type SyncErrorOutputReceiver interface {
	SyncWarningf(string, ...Any)
	SyncErrorf(string, ...Any)
	SyncFatalf(string, ...Any)
}

// ErrorReceiver is an interface
type ErrorReceiver interface {
	Warning(error)
	Error(error)
	Fatal(error)
}

// SyncErrorReceiver is an interface
type SyncErrorReceiver interface {
	SyncWarning(error)
	SyncError(error)
	SyncFatal(error)
}

// SyncLogger is a logger that implements syncronous methods.
type SyncLogger interface {
	Listenable
	SyncTriggerable
	SyncOutputReceiver
	SyncErrorOutputReceiver
	SyncErrorReceiver
}

// AsyncLogger is a logger that implements async methods.
type AsyncLogger interface {
	Listenable
	Triggerable
	OutputReceiver
	ErrorOutputReceiver
	ErrorReceiver
}

// FullLogger is every possible interface.
type FullLogger interface {
	Listenable
	SyncTriggerable
	SyncOutputReceiver
	SyncErrorOutputReceiver
	SyncErrorReceiver
	Triggerable
	OutputReceiver
	ErrorOutputReceiver
	ErrorReceiver
}
