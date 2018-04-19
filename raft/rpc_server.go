package raft

import (
	"net"
	"net/rpc"

	"github.com/blend/go-sdk/exception"
	"github.com/blend/go-sdk/logger"
	"github.com/blend/go-sdk/worker"
)

// NewServer returns a new server.
func NewServer() *Server {
	return &Server{
		methods: &ServerMethods{},
		latch:   &worker.Latch{},
	}
}

// NewServerFromConfig creates a new raft node from a config.
func NewServerFromConfig(cfg *Config) *Server {
	return NewServer().WithBindAddr(cfg.GetBindAddr())
}

// Server is the base entity / fsm manager for the raft election process.
type Server struct {
	bindAddr string
	log      *logger.Logger
	latch    *worker.Latch

	server   *rpc.Server
	listener *net.TCPListener

	methods *ServerMethods
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

// WithBindAddr sets the bind address.
func (s *Server) WithBindAddr(bindAddr string) *Server {
	s.bindAddr = bindAddr
	return s
}

// BindAddr returns the bind address for the rpc server.
func (s *Server) BindAddr() string {
	return s.bindAddr
}

// SetAppendEntriesHandler sets the append entries handler.
func (s *Server) SetAppendEntriesHandler(handler func(*AppendEntries, *AppendEntriesResults) error) {
	s.methods.appendEntriesHandler = handler
}

// SetRequestvoteHandler sets the request vote handler.
func (s *Server) SetRequestvoteHandler(handler func(*RequestVote, *RequestVoteResults) error) {
	s.methods.requestVoteHandler = handler
}

// Start starts the server.
func (s *Server) Start() (err error) {
	s.latch.Starting()
	var addr *net.TCPAddr
	addr, err = net.ResolveTCPAddr("tcp", s.bindAddr)
	if err != nil {
		err = exception.Wrap(err)
		return
	}

	s.listener, err = net.ListenTCP("tcp", addr)
	if err != nil {
		err = exception.Wrap(err)
		return
	}

	s.server = rpc.NewServer()
	err = s.server.Register(s.methods)
	if err != nil {
		err = exception.Wrap(err)
		return
	}

	go func() {
		defer s.latch.Stopped()
		s.latch.Started()
		s.server.Accept(s.listener)
	}()
	<-s.latch.NotifyStarted()
	return
}

// Close closes the server.
func (s *Server) Close() error {
	s.latch.Stopped()
	return exception.Wrap(s.listener.Close())
}

// ServerMethods are the methods we register with the rpc server.
type ServerMethods struct {
	appendEntriesHandler func(*AppendEntries, *AppendEntriesResults) error
	requestVoteHandler   func(*RequestVote, *RequestVoteResults) error
}

// AppendEntries calls the append entries handler.
func (sm *ServerMethods) AppendEntries(args *AppendEntries, res *AppendEntriesResults) error {
	if sm.appendEntriesHandler == nil {
		return nil
	}
	return sm.appendEntriesHandler(args, res)
}

// RequestVote calls the request vote handler.
func (sm *ServerMethods) RequestVote(args *RequestVote, res *RequestVoteResults) error {
	if sm.requestVoteHandler == nil {
		return nil
	}
	return sm.requestVoteHandler(args, res)
}
