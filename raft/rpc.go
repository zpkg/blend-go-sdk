package raft

// Vote is a command to request a leader vote.
type Vote struct {
	Term      uint64
	Candidate string
}

// VoteResponse is the response from nodes during an election.
type VoteResponse struct {
	Term    uint64
	Granted bool
}

// Heartbeat is sent by leaders to prevent election timeouts.
type Heartbeat struct {
	Leader string
	Term   uint64
}

// HeartbeatResponse is a response to a heartbeat.
type HeartbeatResponse struct {
	Term    uint64
	Success bool
}
