package raft

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/blend/go-sdk/exception"
	"github.com/blend/go-sdk/logger"
)

var (
	_ Server = &RPCServer{}
)

// NewRPCServer returns a new roc server.
func NewRPCServer() *RPCServer {
	return &RPCServer{}
}

// RPCServer is the net/rpc implementation of the raft server components.
type RPCServer struct {
	bindAddr      string
	timeout       time.Duration
	requestVote   RequestVoteHandler
	appendEntries AppendEntriesHandler
	log           *logger.Logger
	server        *http.Server
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

// SetAppendEntriesHandler should register the append entries handler.
func (s *RPCServer) SetAppendEntriesHandler(handler AppendEntriesHandler) {
	s.appendEntries = handler
}

// AppendEntriesHandler returns the append entries handler.
func (s *RPCServer) AppendEntriesHandler() AppendEntriesHandler {
	return s.appendEntries
}

// SetRequestVoteHandler should register the request vote handler.
func (s *RPCServer) SetRequestVoteHandler(handler RequestVoteHandler) {
	s.requestVote = handler
}

// RequestVoteHandler returns the request vote handler.
func (s *RPCServer) RequestVoteHandler() RequestVoteHandler {
	return s.requestVote
}

func (s *RPCServer) handle(action http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				if s.log != nil {
					s.log.Fatal(exception.New(r))
				}
			}
		}()
		if s.log != nil {
			s.log.Trigger(logger.NewWebRequestStartEvent(req).WithFlag(logger.Flag("rpc.handler.start")))
			start := time.Now()
			instrumented := logger.NewResponseWriter(w)
			defer func() {
				s.log.Trigger(logger.NewWebRequestEvent(req).WithFlag(logger.Flag("rpc.handler.complete")).
					WithStatusCode(instrumented.StatusCode()).
					WithElapsed(time.Since(start)).
					WithContentLength(int64(instrumented.ContentLength())))
			}()
			action(instrumented, req)
			return
		}
		action(w, req)
	}
}

func (s *RPCServer) appendEntriesHandler(w http.ResponseWriter, req *http.Request) {
	var args AppendEntries
	if err := s.decode(&args, req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var res AppendEntriesResults
	s.appendEntries(&args, &res)
	w.WriteHeader(http.StatusOK)
	if err := s.encode(res, w); err != nil {
		if s.log != nil {
			s.log.Error(err)
		}
	}
}

func (s *RPCServer) requestVoteHandler(w http.ResponseWriter, req *http.Request) {
	var args RequestVote
	if err := s.decode(&args, req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var res RequestVoteResults
	s.requestVote(&args, &res)
	w.WriteHeader(http.StatusOK)
	if err := s.encode(res, w); err != nil {
		if s.log != nil {
			s.log.Error(err)
		}
	}
}

func (s *RPCServer) decode(obj interface{}, req *http.Request) error {
	if req.Body == nil {
		return exception.New("request body unset")
	}
	defer req.Body.Close()
	return exception.New(json.NewDecoder(req.Body).Decode(obj))
}

func (s *RPCServer) encode(obj interface{}, w http.ResponseWriter) error {
	return exception.New(json.NewEncoder(w).Encode(obj))
}

// Start starts the server.
func (s *RPCServer) Start() error {
	mux := http.NewServeMux()
	mux.HandleFunc(fmt.Sprintf("/%s", RPCMethodRequestVote), s.handle(s.requestVoteHandler))
	mux.HandleFunc(fmt.Sprintf("/%s", RPCMethodAppendEntries), s.handle(s.appendEntriesHandler))

	s.server = &http.Server{
		Addr:         s.bindAddr,
		ReadTimeout:  s.timeout,
		WriteTimeout: s.timeout,
		Handler:      mux,
	}

	go s.server.ListenAndServe()
	return nil
}

// Stop stops the server.
func (s *RPCServer) Stop() error {
	timeout, cancel := context.WithTimeout(context.TODO(), time.Second)
	defer cancel()
	return exception.Wrap(s.server.Shutdown(timeout))
}
