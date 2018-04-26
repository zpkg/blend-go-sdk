package raft

// RequestVote is a command to request a leader vote.
type RequestVote struct {
	ID           string
	Term         uint64
	LastLogIndex uint64
	LastLogTerm  uint64
}

// RequestVoteResults is the response from nodes during an election.
type RequestVoteResults struct {
	ID      string
	Term    uint64
	Granted bool
}
