package raft

// ServerMethods are the methods we register with the rpc server.
type ServerMethods struct {
	appendEntriesHandler AppendEntriesHandler
	requestVoteHandler   RequestVoteHandler
}

// SetAppendEntriesHandler sets the append entries handler.
func (sm *ServerMethods) SetAppendEntriesHandler(handler AppendEntriesHandler) {
	sm.appendEntriesHandler = handler
}

// SetRequestVoteHandler sets the request vote handler.
func (sm *ServerMethods) SetRequestVoteHandler(handler RequestVoteHandler) {
	sm.requestVoteHandler = handler
}

// AppendEntries calls the append entries handler.
func (sm ServerMethods) AppendEntries(args *AppendEntries, res *AppendEntriesResults) error {
	if sm.appendEntriesHandler == nil {
		return nil
	}
	return sm.appendEntriesHandler(args, res)
}

// RequestVote calls the request vote handler.
func (sm ServerMethods) RequestVote(args *RequestVote, res *RequestVoteResults) error {
	if sm.requestVoteHandler == nil {
		return nil
	}
	return sm.requestVoteHandler(args, res)
}

// AppendEntriesHandler returns the append entries handler.
func (sm *ServerMethods) AppendEntriesHandler() AppendEntriesHandler { return sm.appendEntriesHandler }

// RequestVoteHandler returns the request vote handler.
func (sm *ServerMethods) RequestVoteHandler() RequestVoteHandler { return sm.requestVoteHandler }
