package jobkit

// JobStatus is a status type for jobs.
type JobStatus string

// Job Statuses
const (
	JobStatusStarted   JobStatus = "started"
	JobStatusCancelled JobStatus = "cancelled"
	JobStatusComplete  JobStatus = "complete"
	JobStatusFailed    JobStatus = "failed"
	JobStatusBroken    JobStatus = "broken"
	JobStatusFixed     JobStatus = "fixed"
)
