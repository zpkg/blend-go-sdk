package raft

import "fmt"

var (
	// Asserts that mock transport implements transport.
	_ Transport = &MockTransport{}
)

// NewMockTransport returns a new mock transport.
func NewMockTransport(remoteAddr string) *MockTransport {
	return &MockTransport{
		remoteAddr: remoteAddr,
	}
}

// MockTransport is a transport used for testing.
type MockTransport struct {
	remoteAddr         string
	requestVoteHandler func(*Vote) (*VoteResponse, error)
	heartbeatHandler   func(*Heartbeat) (*HeartbeatResponse, error)
}

// RemoteAddr retruns the remote addr.
func (mt *MockTransport) RemoteAddr() string {
	return mt.remoteAddr
}

// SetRequestVoteHandler sets the request vote handler.
func (mt *MockTransport) SetRequestVoteHandler(handler func(*Vote) (*VoteResponse, error)) {
	mt.requestVoteHandler = handler
}

// RequestVote implements the request vote handler.
func (mt *MockTransport) RequestVote(args *Vote) (*VoteResponse, error) {
	if mt.requestVoteHandler == nil {
		return nil, fmt.Errorf("request vote handler unset")
	}
	return mt.requestVoteHandler(args)
}

// SetHeartbeatHandler sets the heartbeat handler.
func (mt *MockTransport) SetHeartbeatHandler(handler func(*Heartbeat) (*HeartbeatResponse, error)) {
	mt.heartbeatHandler = handler
}

// Heartbeat implements the heartbeat handler.
func (mt *MockTransport) Heartbeat(rpc *Heartbeat) (*HeartbeatResponse, error) {
	if mt.heartbeatHandler == nil {
		return nil, fmt.Errorf("heartbeat handler unset")
	}
	return mt.heartbeatHandler(rpc)
}

// Close implements the mock transport closer.
func (mt *MockTransport) Close() error { return nil }
