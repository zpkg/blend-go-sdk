package logger

import "os"

// FatalExit creates a basic logger and calls `SyncFatalExit` on it.
func FatalExit(err error) {
	if err == nil {
		return
	}
	log := Sync()
	log.Enable(Fatal)
	log.Fatal(err)
	os.Exit(1)
}
