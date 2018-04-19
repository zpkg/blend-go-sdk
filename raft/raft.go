package raft

import (
	"sync"
	"sync/atomic"
	"time"

	"github.com/blend/go-sdk/logger"
	"github.com/blend/go-sdk/worker"
)

const (
	// DefaultLeaderCheckTick is the tick rate for the leader check.
	DefaultLeaderCheckTick = 250 * time.Millisecond

	// DefaultHeartbeatTick is the tick rate for leaders to send heartbeats.
	DefaultHeartbeatTick = 500 * time.Millisecond
)

// FSMState is a raft fsm state.
type FSMState int32

func (fsm FSMState) String() string {
	switch fsm {
	case FSMStateUnset:
		return "unset"
	case FSMStateFollower:
		return "follower"
	case FSMStateCandidate:
		return "candidate"
	case FSMStateLeader:
		return "leader"
	default:
		return "unknown"
	}
}

const (
	// FSMStateUnset is the unset fsm state.
	FSMStateUnset FSMState = 0
	// FSMStateFollower is the follower state.
	FSMStateFollower FSMState = 1
	// FSMStateCandidate is the follower state.
	FSMStateCandidate FSMState = 2
	// FSMStateLeader is the follower state.
	FSMStateLeader FSMState = 3
)

// NewFromConfig creates a new raft node from a config.
func NewFromConfig(cfg *Config) *Raft {
	return &Raft{
		id:  cfg.GetIdentifier(),
		cfg: cfg,
	}
}

// Raft represents a raft node and all the state machine
// componentry required.
type Raft struct {
	id  string
	cfg *Config
	log *logger.Logger

	currentTerm       uint64
	votedFor          string
	lastLeaderContact time.Time
	lastVoteGranted   time.Time

	// state is the current fsm state
	state int32

	// these are unused
	commitIndex uint64            // the highest log entry known to be committed
	lastApplied uint64            // index of highest log entry applied to state machine
	nextIndex   map[string]uint64 // for each peer, the index of the next log entry to send to that peer
	matchIndex  map[string]uint64 // for each peer, index of highest log entry known to be replicated on each peer

	server *Server
	peers  []*Client

	leaderCheckTick *worker.Interval
	heartbeatTick   *worker.Interval

	leaderHandler    func()
	candidateHandler func()
	followerHandler  func()
}

// ID is the raft node identifier.
func (r *Raft) ID() string {
	return r.id
}

// Config returns the config.
func (r *Raft) Config() *Config {
	return r.cfg
}

// State returns the current raft state.
func (r *Raft) State() FSMState {
	return FSMState(atomic.LoadInt32(&r.state))
}

// SetLeaderHandler sets the leader handler.
func (r *Raft) SetLeaderHandler(handler func()) {
	r.leaderHandler = handler
}

// SetCandidateHandler sets the leader handler.
func (r *Raft) SetCandidateHandler(handler func()) {
	r.candidateHandler = handler
}

// SetFollowerHandler sets the leader handler.
func (r *Raft) SetFollowerHandler(handler func()) {
	r.followerHandler = handler
}

// WithLogger sets the logger.
func (r *Raft) WithLogger(log *logger.Logger) *Raft {
	r.log = log
	return r
}

// Logger returns the logger.
func (r *Raft) Logger() *logger.Logger {
	return r.log
}

// WithPeer adds a peer.
func (r *Raft) WithPeer(peer *Client) *Raft {
	r.peers = append(r.peers, peer)
	return r
}

// Peers returns the raft peers.
func (r *Raft) Peers() []*Client {
	return r.peers
}

// Start starts the raft node.
func (r *Raft) Start() error {
	r.infof("node starting")
	defer func() {
		r.infof("node started")
	}()

	r.server = NewServerFromConfig(r.cfg).WithLogger(r.log)
	r.server.SetAppendEntriesHandler(r.appendEntriesHandler)
	r.server.SetRequestvoteHandler(r.requestVoteHandler)

	r.infof("node rpc server starting, listening on: %s", r.cfg.GetBindAddr())
	err := r.server.Start()
	if err != nil {
		return err
	}

	r.infof("node beginning internal tickers")
	r.leaderCheckTick = worker.NewInterval(r.leaderCheck, DefaultLeaderCheckTick).WithDelay(r.cfg.GetElectionTimeout())
	r.heartbeatTick = worker.NewInterval(r.heartbeat, DefaultHeartbeatTick).WithDelay(r.cfg.GetElectionTimeout())
	r.leaderCheckTick.Start()
	r.heartbeatTick.Start()
	return nil
}

// Stop stops the node.
func (r *Raft) Stop() error {
	if r.leaderCheckTick != nil {
		r.leaderCheckTick.Stop()
	}
	if r.heartbeatTick != nil {
		r.heartbeatTick.Stop()
	}

	return r.server.Close()
}

// utility methods.

func (r *Raft) transitionTo(state FSMState) {
	r.infof("transitioning to %s", state)

	switch state {
	case FSMStateFollower:
		atomic.StoreInt32(&r.state, int32(FSMStateFollower))
		if r.followerHandler != nil {
			r.followerHandler()
		}
	case FSMStateCandidate:
		atomic.StoreInt32(&r.state, int32(FSMStateCandidate))
		if r.candidateHandler != nil {
			r.candidateHandler()
		}
	case FSMStateLeader:
		atomic.StoreInt32(&r.state, int32(FSMStateLeader))
		if r.leaderHandler != nil {
			r.leaderHandler()
		}
	}
}

// LeaderCheck is the action that fires on a heartbeat to check if the leader lease has expired.
func (r *Raft) leaderCheck() error {

	if r.State() == FSMStateLeader {
		return nil
	}

	now := time.Now().UTC()
	if r.lastLeaderContact.IsZero() || now.Sub(r.lastLeaderContact) > r.cfg.GetLeaderLeaseTimeout() {
		if r.lastVoteGranted.IsZero() || now.Sub(r.lastVoteGranted) > r.cfg.GetElectionTimeout() {
			r.infof("leader check fails")
			return r.election()
		}
	}

	return nil
}

func (r *Raft) election() error {
	r.transitionTo(FSMStateCandidate)
	return r.requestVote()
}

func (r *Raft) requestVote() error {
	r.votedFor = r.id
	r.currentTerm = r.currentTerm + 1

	voteRequest := RequestVote{
		Term:      r.currentTerm,
		Candidate: r.id,
	}

	results := make(chan *RequestVoteResults, len(r.peers))
	errs := make(chan error, len(r.peers))
	wg := sync.WaitGroup{}
	wg.Add(len(r.peers))
	for _, peer := range r.peers {
		go func(c *Client) {
			defer wg.Done()

			r.infof("requesting vote from %s", c.RemoteAddr())
			res, err := c.RequestVote(&voteRequest)
			if err != nil {
				errs <- err
			} else {
				r.infof("got result from %s: %v", c.RemoteAddr(), res.Granted)
				results <- res
			}
		}(peer)
	}
	wg.Wait()

	r.infof("election complete")
	if r.countVotes(results) { // we're now the leader.
		r.transitionTo(FSMStateLeader)
	} else {
		r.transitionTo(FSMStateFollower)
	}

	return nil
}

func (r *Raft) countVotes(results chan *RequestVoteResults) bool {
	// tabulate results
	total := len(results)
	votesFor := 0

	if required == 0 {
		return true
	}

	for index := 0; index < required; index++ {
		result := <-results

		if result.Granted {
			votesFor = votesFor + 1
		}
	}

	if total == 1 && votesFor == 1 {
		return true
	}
	if total == 2 && votesFor >= 1 {
		return true
	}
	if total == 3 && votesFor >= 2 {
		return true
	}
	if total == 4 && votesFor >= 2 {
		return true
	}

	return votesFor > total>>1
}

func (r *Raft) heartbeat() error {
	if r.State() != FSMStateLeader {
		return nil
	}

	args := AppendEntries{
		Term:     r.currentTerm,
		LeaderID: r.votedFor,
	}

	results := make(chan *AppendEntriesResults, len(r.peers))
	errs := make(chan error, len(r.peers))
	wg := sync.WaitGroup{}
	wg.Add(len(r.peers))
	for _, peer := range r.peers {
		go func(c *Client) {
			defer wg.Done()

			res, err := c.AppendEntries(&args)
			if err != nil {
				errs <- err
			} else {
				results <- res
			}

		}(peer)
	}
	wg.Wait()

	/*
		for index := 0; index < len(errs); index++ {
			r.err(<-errs)
		}
	*/

	return nil
}

func (r *Raft) appendEntriesHandler(args *AppendEntries, res *AppendEntriesResults) error {
	// Reply false if term < currentTerm
	if args.Term < r.currentTerm {
		res = &AppendEntriesResults{
			Success: false,
			Term:    r.currentTerm,
		}
		return nil
	}

	r.votedFor = args.LeaderID
	r.lastLeaderContact = time.Now().UTC()
	res = &AppendEntriesResults{
		Success: true,
		Term:    r.currentTerm,
	}

	return nil
}

func (r *Raft) requestVoteHandler(args *RequestVote, res *RequestVoteResults) error {
	if args.Term < r.currentTerm {
		r.infof("requestVote from %s fails because args.Term(%d) < currentTerm(%d)", args.Candidate, args.Term, r.currentTerm)
		res = &RequestVoteResults{
			Term:    r.currentTerm,
			Granted: false,
		}
		return nil
	}

	if len(r.votedFor) == 0 || r.votedFor == args.Candidate {
		r.lastVoteGranted = time.Now().UTC()
		r.votedFor = args.Candidate
		res = &RequestVoteResults{
			Term:    args.Term,
			Granted: true,
		}
		return nil
	}

	r.infof("requestVote from %s fails because r.votedFor == %s", args.Candidate, r.votedFor)
	res = &RequestVoteResults{
		Term:    args.Term,
		Granted: false,
	}
	return nil
}

func (r *Raft) dialPeers() error {
	for _, peer := range r.peers {
		r.infof("node dialing peer: %s", peer.RemoteAddr())
		if err := peer.Open(); err != nil {
			return err
		}
	}
	return nil
}

func (r *Raft) infof(format string, args ...interface{}) {
	if r.log != nil {
		r.log.SubContext("raft").SubContext(r.id).SyncInfof(format, args...)
	}
}

func (r *Raft) err(err error) {
	if r.log != nil {
		r.log.SubContext("raft").SubContext(r.id).SyncError(err)
	}
}
