package logger

import "context"

// MaybeTrigger triggers an event if the logger is set.
func MaybeTrigger(log Triggerable, e Event) {
	if log == nil {
		return
	}
	log.Trigger(context.Background(), e)
}

// MaybeInfof triggers Infof if the logger is set.
func MaybeInfof(log InfofReceiver, format string, args ...interface{}) {
	if log == nil {
		return
	}
	log.Infof(format, args...)
}

// MaybeDebugf triggers Debugf if the logger is set.
func MaybeDebugf(log DebugfReceiver, format string, args ...interface{}) {
	if log == nil {
		return
	}
	log.Debugf(format, args...)
}

// MaybeWarningf triggers Warningf if the logger is set.
func MaybeWarningf(log WarningfReceiver, format string, args ...interface{}) {
	if log == nil {
		return
	}
	log.Warningf(format, args...)
}

// MaybeWarning triggers Warning if the logger is set.
func MaybeWarning(log WarningReceiver, err error) {
	if log == nil || err == nil {
		return
	}
	log.Warning(err)
}

// MaybeErrorf triggers Errorf if the logger is set.
func MaybeErrorf(log ErrorfReceiver, format string, args ...interface{}) {
	if log == nil {
		return
	}
	log.Errorf(format, args...)
}

// MaybeError triggers Error if the logger is set.
func MaybeError(log ErrorReceiver, err error) {
	if log == nil || err == nil {
		return
	}
	log.Error(err)
}

// MaybeFatalf triggers Fatalf if the logger is set.
func MaybeFatalf(log FatalfReceiver, format string, args ...interface{}) {
	if log == nil {
		return
	}
	log.Fatalf(format, args...)
}

// MaybeFatal triggers Fatal if the logger is set.
func MaybeFatal(log FatalReceiver, err error) {
	if log == nil || err == nil {
		return
	}
	log.Fatal(err)
}
