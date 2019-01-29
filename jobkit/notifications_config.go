package jobkit

import (
	"github.com/blend/go-sdk/configutil"
)

// NotificationsConfig is something you can use to give your jobs some knobs to turn
// from configuration.
// You can use this job config by embedding it into your larger job config struct.
type NotificationsConfig struct {
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

// NotifyOnStartOrDefault returns a value or a default.
func (nc NotificationsConfig) NotifyOnStartOrDefault() bool {
	return configutil.CoalesceBool(nc.NotifyOnStart, false)
}

// NotifyOnSuccessOrDefault returns a value or a default.
func (nc NotificationsConfig) NotifyOnSuccessOrDefault() bool {
	return configutil.CoalesceBool(nc.NotifyOnSuccess, false)
}

// NotifyOnFailureOrDefault returns a value or a default.
func (nc NotificationsConfig) NotifyOnFailureOrDefault() bool {
	return configutil.CoalesceBool(nc.NotifyOnFailure, false)
}

// NotifyOnBrokenOrDefault returns a value or a default.
func (nc NotificationsConfig) NotifyOnBrokenOrDefault() bool {
	return configutil.CoalesceBool(nc.NotifyOnBroken, true)
}

// NotifyOnFixedOrDefault returns a value or a default.
func (nc NotificationsConfig) NotifyOnFixedOrDefault() bool {
	return configutil.CoalesceBool(nc.NotifyOnFixed, true)
}
