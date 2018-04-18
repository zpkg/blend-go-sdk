package raft

// Transport is the interface a type must implement
// to allow us to communicate over it.
// It generally is analogous to a connection to a node.
type Transport interface {
	// RemoteAddr is the remote address of the transport.
	RemoteAddr() string

	// RequestVote kicks off an election (a request for votes).
	// `id` is the candidate id, usually the node applying to be leader.
	RequestVote(args *Vote) (*VoteResponse, error)
	// Heartbeat sends a heartbeat to the node.
	Heartbeat(*Heartbeat) (*HeartbeatResponse, error)

	// Close disconnects the transport.
	Close() error
}
