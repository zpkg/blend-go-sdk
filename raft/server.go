package raft

import (
	"net"
	"net/rpc"

	"github.com/blend/go-sdk/exception"
	"github.com/blend/go-sdk/logger"
	"github.com/blend/go-sdk/worker"
)

// New returns a new server.
func New() *Server {
	return &Server{}
}

// NewFromConfig creates a new raft node from a config.
func NewFromConfig(cfg *Config) *Server {
	var peers []Transport
	for _, peer := range cfg.GetPeers() {
		peers = append(peers, NewRPCTransport(peer))
	}
	return &Server{
		identifier: cfg.GetIdentifier(),
		bindAddr:   cfg.GetBindAddr(),
		peers:      peers,
	}
}

// Server is the base entity / fsm manager for the raft election process.
type Server struct {
	identifier string
	bindAddr   string

	leader      string
	currentTerm uint64

	server   *rpc.Server
	listener *net.TCPListener

	peers []Transport

	heartbeatSender *worker.Interval
	leaderCheck     *worker.Interval

	log *logger.Logger
	cfg *Config
}

// WithLogger sets the logger.
func (s *Server) WithLogger(log *logger.Logger) *Server {
	s.log = log
	return s
}

// Logger returns the logger.
func (s *Server) Logger() *logger.Logger {
	return s.log
}

// WithID sets the identifier.
func (s *Server) WithID(id string) *Server {
	s.identifier = id
	return s
}

// ID returns the identifier.
func (s *Server) ID() string {
	return s.identifier
}

// IsLeader returns if a server is a leader.
func (s *Server) IsLeader() bool {
	return s.leader == s.identifier
}

// WithPeer adds a peer to the server.
func (s *Server) WithPeer(peer Transport) *Server {
	s.peers = append(s.peers, peer)
	return s
}

// Initialize initializes the node.
func (s *Server) Initialize() (err error) {
	s.log.SyncDebugf("%s raft initializing, listening on: %s", s.identifier, s.bindAddr)

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

	for _, peer := range s.peers {
		err = peer.Open()
		if err != nil {
			err = exception.Wrap(err)
			return
		}
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
	if s.leader == rpc.Leader {
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
