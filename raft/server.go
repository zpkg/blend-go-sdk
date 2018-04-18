package raft

import (
	"net"
	"net/rpc"
	"time"

	"github.com/blend/go-sdk/exception"
	"github.com/blend/go-sdk/logger"
	"github.com/blend/go-sdk/worker"
)

// Server is the base entity / fsm manager for the raft election process.
type Server struct {
	identifier  string
	bindAddr    string
	lastContact time.Time

	votedFor    string
	currentTerm uint64

	server   *rpc.Server
	listener *net.TCPListener

	heartbeatSender *worker.Interval

	peers []Transport

	log *logger.Logger
	cfg *Config
}

// IsLeader returns if a server is a leader.
func (s *Server) IsLeader() bool {
	return s.votedFor == s.identifier
}

// WithPeer adds a peer to the server.
func (s *Server) WithPeer(peer Transport) *Server {
	s.peers = append(s.peers, peer)
	return s
}

// Initialize initializes the node.
func (s *Server) Initialize() (err error) {
	var addr *net.TCPAddr
	addr, err = net.ResolveTCPAddr("tcp", s.bindAddr)
	if err != nil {
		err = exception.Wrap(err)
	}

	s.listener, err = net.ListenTCP("tcp", addr)

	s.server = rpc.NewServer()
	err = s.server.Register(s)
	if err != nil {
		err = exception.Wrap(err)
	}

	go func() { s.server.Accept(s.listener) }()
	return nil
}

// Close closes the server.
func (s *Server) Close() error {
	var err error
	for _, peer := range s.peers {
		err = peer.Close()
		if err != nil {
			return exception.Wrap(err)
		}
	}
	return exception.Wrap(s.listener.Close())
}

// RequestVote handles a vote rpc.
func (s *Server) RequestVote(rpc *Vote) (*VoteResponse, error) {
	if s.currentTerm < rpc.Term {
		return &VoteResponse{
			Term:    s.currentTerm,
			Granted: false,
		}, nil
	}

	return &VoteResponse{
		Term:    s.currentTerm,
		Granted: true,
	}, nil
}

// Heartbeat handles a heartbeat rpc.
func (s *Server) Heartbeat(rpc *Heartbeat) (*HeartbeatResponse, error) {
	if s.votedFor == rpc.Leader {
		s.currentTerm = rpc.Term
		return &HeartbeatResponse{
			Success: true,
			Term:    s.currentTerm,
		}, nil
	}

	// update last contact information
	return &HeartbeatResponse{
		Success: false,
		Term:    s.currentTerm,
	}, nil
}
