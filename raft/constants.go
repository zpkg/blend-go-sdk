package raft

import "time"

const (
	// RPCMethodRequestVote is an rpc method.
	RPCMethodRequestVote = "ServerMethods.RequestVote"
	// RPCMethodAppendEntries is an rpc method.
	RPCMethodAppendEntries = "ServerMethods.AppendEntries"

	// DefaultClientTimeout is a default.
	DefaultClientTimeout = 500 * time.Millisecond
)
