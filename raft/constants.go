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
