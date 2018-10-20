package cron

import "github.com/blend/go-sdk/exception"

const (
	// ErrJobNotLoaded is a common error.
	ErrJobNotLoaded exception.Class = "job not loaded"

	// ErrJobAlreadyLoaded is a common error.
	ErrJobAlreadyLoaded exception.Class = "job already loaded"

	// ErrTaskNotFound is a common error.
	ErrTaskNotFound exception.Class = "task not found"

	// ErrTaskCancelled is a common error.
	ErrTaskCancelled exception.Class = "task cancelled"
)

// IsJobNotLoaded returns if the error is a job not loaded error.
func IsJobNotLoaded(err error) bool {
	return exception.Is(err, ErrJobNotLoaded)
}

// IsJobAlreadyLoaded returns if the error is a job already loaded error.
func IsJobAlreadyLoaded(err error) bool {
	return exception.Is(err, ErrJobAlreadyLoaded)
}

// IsTaskNotFound returns if the error is a task not found error.
func IsTaskNotFound(err error) bool {
	return exception.Is(err, ErrTaskNotFound)
}

// IsTaskCancelled returns if the error is a task not found error.
func IsTaskCancelled(err error) bool {
	return exception.Is(err, ErrTaskCancelled)
}
