package logger

import "os"

// MaybeFatalExit will print the error and exit the process
// with exit(1) if the error isn't nil.
func MaybeFatalExit(err error) {
	if err == nil {
		return
	}
	FatalExit(err)
}

// FatalExit will print the error and exit the process with exit(1).
func FatalExit(err error) {
	MustNew(OptOutput(os.Stderr)).Fatal(err)
	os.Exit(1)
}
