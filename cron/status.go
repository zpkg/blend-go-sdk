package cron

// Status is a status object
type Status struct {
	Jobs  []JobMeta
	Tasks map[string]TaskInvocation
}
