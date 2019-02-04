package jobkit

import (
	"github.com/blend/go-sdk/configutil"
	"github.com/blend/go-sdk/uuid"
)

// JobConfig is something you can use to give your jobs some knobs to turn
// from configuration.
// You can use this job config by embedding it into your larger job config struct.
type JobConfig struct {
	// Name is the name of the job.
	Name string `json:"name" yaml:"name"`
	// Schedule returns the job schedule.
	Schedule string `json:"schedule" yaml:"schedule"`

	// NotifyOnStart governs if we should send notifications job start.
	NotifyOnStart *bool `json:"notifyOnStart" yaml:"notifyOnStart"`
	// NotifyOnSuccess governs if we should send notifications on any success.
	NotifyOnSuccess *bool `json:"notifyOnSuccess" yaml:"notifyOnSuccess"`
	// NotifyOnFailure governs if we should send notifications on any failure.
	NotifyOnFailure *bool `json:"notifyOnFailure" yaml:"notifyOnFailure"`
	// NotifyOnBroken governs if we should send notifications on a success => failure transition.
	NotifyOnBroken *bool `json:"notifyOnBroken" yaml:"notifyOnBroken"`
	// NotifyOnFixed governs if we should send notifications on a failure => success transition.
	NotifyOnFixed *bool `json:"notifyOnFixed" yaml:"notifyOnFixed"`
}

// NameOrDefault returns the job name or a default (uuid v4).
func (jc JobConfig) NameOrDefault() string {
	return configutil.CoalesceString(jc.Name, uuid.V4().String())
}

// ScheduleOrDefault returns the schedule or a default (every 5 minutes).
func (jc JobConfig) ScheduleOrDefault() string {
	return configutil.CoalesceString(jc.Schedule, "* */5 * * * * *")
}

// NotifyOnStartOrDefault returns a value or a default.
func (jc JobConfig) NotifyOnStartOrDefault() bool {
	return configutil.CoalesceBool(jc.NotifyOnStart, false)
}

// NotifyOnSuccessOrDefault returns a value or a default.
func (jc JobConfig) NotifyOnSuccessOrDefault() bool {
	return configutil.CoalesceBool(jc.NotifyOnSuccess, false)
}

// NotifyOnFailureOrDefault returns a value or a default.
func (jc JobConfig) NotifyOnFailureOrDefault() bool {
	return configutil.CoalesceBool(jc.NotifyOnFailure, false)
}

// NotifyOnBrokenOrDefault returns a value or a default.
func (jc JobConfig) NotifyOnBrokenOrDefault() bool {
	return configutil.CoalesceBool(jc.NotifyOnBroken, true)
}

// NotifyOnFixedOrDefault returns a value or a default.
func (jc JobConfig) NotifyOnFixedOrDefault() bool {
	return configutil.CoalesceBool(jc.NotifyOnFixed, true)
}
