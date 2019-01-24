package cron

// Status is a status object
type Status struct {
	Jobs    []*JobScheduler           `json:"jobs"`
	Running map[string]*JobInvocation `json:"running,omitempty"`
}
