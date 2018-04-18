package raft

import (
	"net/rpc"

	"github.com/blend/go-sdk/exception"
)

var (
	// Asserts that RPCTransport implements Transport.
	_ Transport = &RPCTransport{}
)

const (
	// RPCMethodRequestVote is an rpc method.
	RPCMethodRequestVote = "RequestVote"
	// RPCMethodHeatbeat is an rpc method.
	RPCMethodHeatbeat = "Heartbeat"
)

// RPCTransport is an rpc transport over the network.
type RPCTransport struct {
	remoteAddr string
	client     *rpc.Client
}

// RemoteAddr returns the remote address.
func (rt *RPCTransport) RemoteAddr() string { return rt.remoteAddr }

// Open opens the connection and starts the receive server.
func (rt *RPCTransport) Open() (err error) {
	rt.client, err = rpc.Dial("TCP", rt.remoteAddr)
	if err != nil {
		err = exception.Wrap(err)
		return
	}

	return
}

// RequestVote implements the request vote handler.
func (rt *RPCTransport) RequestVote(args *Vote) (*VoteResponse, error) {
	var res VoteResponse
	err := rt.client.Call(RPCMethodRequestVote, args, &res)
	if err != nil {
		return nil, exception.Wrap(err)
	}
	return &res, nil
}

// Heartbeat implements the heartbeat handler.
func (rt *RPCTransport) Heartbeat(args *Heartbeat) (*HeartbeatResponse, error) {
	var res HeartbeatResponse
	err := rt.client.Call(RPCMethodHeatbeat, args, &res)
	if err != nil {
		return nil, exception.Wrap(err)
	}
	return &res, nil
}

// Close closes the transport.
func (rt *RPCTransport) Close() error {
	if rt.client == nil {
		return nil
	}
	return exception.Wrap(rt.client.Close())
}
