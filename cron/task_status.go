package cron

// --------------------------------------------------------------------------------
// task status
// --------------------------------------------------------------------------------

// TaskStatus is the basic format of a status of a task.
type TaskStatus struct {
	Name        string `json:"name"`
	State       State  `json:"state"`
	Status      string `json:"status,omitempty"`
	LastRunTime string `json:"last_run_time,omitempty"`
	NextRunTime string `json:"next_run_time,omitempy"`
	RunningFor  string `json:"running_for,omitempty"`
	Serial      bool   `json:"serial_execution,omitempty"`
}
