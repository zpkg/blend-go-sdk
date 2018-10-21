package cron

import (
	"time"
)

// JobMeta is runtime metadata for a job.
type JobMeta struct {
	Name        string    `json:"name"`
	Job         Job       `json:"job"`
	Disabled    bool      `json:"disabled"`
	NextRunTime time.Time `json:"nextRunTime"`

	Schedule Schedule `json:"-"`

	EnabledProvider                func() bool          `json:"-"`
	SerialProvider                 func() bool          `json:"-"`
	TimeoutProvider                func() time.Duration `json:"-"`
	ShouldTriggerListenersProvider func() bool          `json:"-"`
	ShouldWriteOutputProvider      func() bool          `json:"-"`

	Last *JobInvocation `json:"last"`
}
