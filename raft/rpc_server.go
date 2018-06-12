package raft

import (
	"net"
	"net/rpc"

	"github.com/blend/go-sdk/exception"
	"github.com/blend/go-sdk/logger"
	"github.com/blend/go-sdk/worker"
)

var (
	_ Server = &RPCServer{}
)

// NewRPCServer returns a new roc server.
func NewRPCServer() *RPCServer {
	return &RPCServer{
		latch: &worker.Latch{},
	}
}

// RPCServer is the net/rpc implementation of the raft server components.
type RPCServer struct {
	ServerMethods
	bindAddr string
	log      *logger.Logger
	latch    *worker.Latch

	server   *rpc.Server
	listener *net.TCPListener
}

// WithLogger sets the logger.
func (s *RPCServer) WithLogger(log *logger.Logger) *RPCServer {
	s.log = log
	return s
}

// Logger returns the logger.
func (s *RPCServer) Logger() *logger.Logger {
	return s.log
}

// WithBindAddr sets the bind address.
func (s *RPCServer) WithBindAddr(bindAddr string) *RPCServer {
	s.bindAddr = bindAddr
	return s
}

// BindAddr returns the bind address for the rpc server.
func (s *RPCServer) BindAddr() string {
	return s.bindAddr
}

// Start starts the server.
func (s *RPCServer) Start() (err error) {
	s.latch.Starting()
	var addr *net.TCPAddr
	addr, err = net.ResolveTCPAddr("tcp", s.bindAddr)
	if err != nil {
		err = exception.New(err)
		return
	}

	s.listener, err = net.ListenTCP("tcp", addr)
	if err != nil {
		err = exception.New(err)
		return
	}

	s.server = rpc.NewServer()
	err = s.server.Register(s.ServerMethods)
	if err != nil {
		err = exception.New(err)
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

// Stop stops the server.
func (s *RPCServer) Stop() error {
	s.latch.Stop()
	return exception.New(s.listener.Close())
}
