package raft

import (
	"time"

	"github.com/blend/go-sdk/logger"
)

const (
	// RPCMethodRequestVote is an rpc method.
	RPCMethodRequestVote = "ServerMethods.RequestVote"
	// RPCMethodAppendEntries is an rpc method.
	RPCMethodAppendEntries = "ServerMethods.AppendEntries"

	// DefaultClientTimeout is a default.
	DefaultClientTimeout = 500 * time.Millisecond
	// DefaultServerTimeout is a default.
	DefaultServerTimeout = 500 * time.Millisecond
)

const (
	// FlagRPCHandlerStart is a logger flag.
	FlagRPCHandlerStart = logger.Flag("rpc.handler.start")
	// FlagRPCHandler is a logger flag.
	FlagRPCHandler = logger.Flag("rpc.handler")
)

const (
	// DefaultLeaderCheckInterval is the tick rate for the leader check.
	DefaultLeaderCheckInterval = 2000 * time.Millisecond
	// DefaultHeartbeatInterval is the tick rate for leaders to send heartbeats.
	DefaultHeartbeatInterval = 1000 * time.Millisecond
	// DefaultElectionTimeout is a default.
	DefaultElectionTimeout = 5 * DefaultHeartbeatInterval
	// DefaultElectionBackoffTimeout is a default.
	DefaultElectionBackoffTimeout = DefaultElectionTimeout
	// DefaultPeerDialTimeout is the default peer dial timeout.
	DefaultPeerDialTimeout = time.Second

	// DefaultBindAddr is the default rpc server bind address.
	DefaultBindAddr = ":6060"

	// EnvVarIdentifier is an environment variable.
	EnvVarIdentifier = "RAFT_ID"
	// EnvVarBindAddr is an environment variable.
	EnvVarBindAddr = "RAFT_BIND_ADDR"
	// EnvVarPeers is an environment variable.
	EnvVarPeers = "RAFT_PEERS"
	// EnvVarElectionTimeout is an environment variable.
	EnvVarElectionTimeout = "RAFT_ELECTION_TIMEOUT"
	// EnvVarRaftPeerDialTimeout is an environment variable.
	EnvVarRaftPeerDialTimeout = "RAFT_PEER_DIAL_TIMEOUT"
)
