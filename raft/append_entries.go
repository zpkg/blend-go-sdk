package raft

// AppendEntries is a log propagation request.
type AppendEntries struct {
	// Term represents the leader's term
	Term uint64
	// LeaderID so follower can redirect clients.
	LeaderID string
	// PrevLogIndex is the index of the log entry immediately preceding new ones
	PrevLogIndex uint64
	// PrevLogterm is the term of the PrevLogIndex entry
	PrevLogTerm uint64
	// LeaderCommit is the leader's Commit Index
	LeaderCommit uint64
}

// AppendEntriesResults is the response from an append entries request.
type AppendEntriesResults struct {
	// Term is the current term, for leader to update itself
	Term uint64
	// Success is true if follower contained entry matching the PrevLogIndex and PrevLogTerm
	Success bool
}
