package cron

import "sync"

var _default *JobManager
var _defaultLock = &sync.Mutex{}

// Default returns a shared instance of a JobManager.
func Default() *JobManager {
	if _default == nil {
		_defaultLock.Lock()
		defer _defaultLock.Unlock()

		if _default == nil {
			_default = New()
		}
	}
	return _default
}
