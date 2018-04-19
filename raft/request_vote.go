package raft

// RequestVote is a command to request a leader vote.
type RequestVote struct {
	Term         uint64
	Candidate    string
	LastLogIndex uint64
	LastLogTerm  uint64
}

// RequestVoteResults is the response from nodes during an election.
type RequestVoteResults struct {
	Term    uint64
	Granted bool
}
