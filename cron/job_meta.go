package cron

import (
	"time"
)

// Status is a status object
type Status struct {
	Jobs  []JobMeta
	Tasks map[string]TaskMeta
}

// JobMeta is runtime metadata for a job.
type JobMeta struct {
	Name            string      `json:"name"`
	Job             Job         `json:"-"`
	Disabled        bool        `json:"disabled"`
	Schedule        Schedule    `json:"-"`
	EnabledProvider func() bool `json:"-"`
	NextRunTime     time.Time   `json:"nextRunTime"`
	LastRunTime     time.Time   `json:"lastRunTime"`
}
