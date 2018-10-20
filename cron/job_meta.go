package cron

import (
	"time"
)

// JobMeta is runtime metadata for a job.
type JobMeta struct {
	Name            string      `json:"name"`
	Job             Job         `json:"-"`
	Disabled        bool        `json:"disabled"`
	Schedule        Schedule    `json:"-"`
	EnabledProvider func() bool `json:"-"`
	NextRunTime     time.Time   `json:"nextRunTime"`

	Last *TaskInvocation `json:"last"`
}
