package raft

const (
	// ErrLeader is returned when an op can't be completed on the current leader node.
	ErrLeader Error = "node is the leader"

	// ErrNotLeader is returned when an op can't be completed on a non-leader node.
	ErrNotLeader Error = "node is not the leader"
)
