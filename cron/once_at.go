package cron

import "time"

// OnceAt returns a schedule that fires once at a given time.
// It will never fire again unless reloaded.
func OnceAt(t time.Time) Schedule {
	return OnceAtSchedule{Time: t}
}

// OnceAtSchedule is a schedule.
type OnceAtSchedule struct {
	Time time.Time
}

// GetNextRunTime returns the next runtime.
func (oa OnceAtSchedule) GetNextRunTime(after *time.Time) *time.Time {
	if after == nil {
		return &oa.Time
	}
	if oa.Time.After(*after) {
		return &oa.Time
	}
	return nil
}
