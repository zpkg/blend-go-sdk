package raft

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/blend/go-sdk/logger"
	"github.com/blend/go-sdk/util"
	"github.com/blend/go-sdk/uuid"
	"github.com/blend/go-sdk/worker"
)

// New creates a new empty raft node.
func New() *Raft {
	return &Raft{
		id:       uuid.V4().String(),
		state:    int32(FSMStateFollower),
		bindAddr: DefaultBindAddr,
	}
}

// NewFromConfig creates a new raft node from a config.
func NewFromConfig(cfg *Config) *Raft {
	return New().
		WithID(cfg.GetIdentifier()).
		WithBindAddr(cfg.GetBindAddr()).
		WithLeaderLeaseTimeout(cfg.GetLeaderLeaseTimeout()).
		WithElectionTimeout(cfg.GetElectionTimeout())
}

// Raft represents a raft node and all the state machine
// componentry required.
type Raft struct {
	id                 string
	log                *logger.Logger
	bindAddr           string
	leaderLeaseTimeout time.Duration
	electionTimeout    time.Duration

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

// WithID sets the identifier for the node.
func (r *Raft) WithID(id string) *Raft {
	r.id = id
	return r
}

// ID is the raft node identifier.
func (r *Raft) ID() string {
	return r.id
}

// WithBindAddr sets the rpc server bind address.
func (r *Raft) WithBindAddr(bindAddr string) *Raft {
	r.bindAddr = bindAddr
	return r
}

// BindAddr returns the rpc server bind address.
func (r *Raft) BindAddr() string {
	return util.Coalesce.String(r.bindAddr, DefaultBindAddr)
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

// WithLeaderLeaseTimeout sets the leader lease timeout.
func (r *Raft) WithLeaderLeaseTimeout(d time.Duration) *Raft {
	r.leaderLeaseTimeout = d
	return r
}

// RandomLeaderLeaseTimeout returns a random leader lease timeout.
func (r *Raft) RandomLeaderLeaseTimeout() time.Duration {
	return randomTimeout(util.Coalesce.Duration(r.leaderLeaseTimeout, DefaultLeaderLeaseTimeout))
}

// WithElectionTimeout sets the election timeout.
func (r *Raft) WithElectionTimeout(d time.Duration) *Raft {
	r.electionTimeout = d
	return r
}

// RandomElectionTimeout returns a random election timeout.
func (r *Raft) RandomElectionTimeout() time.Duration {
	return randomTimeout(util.Coalesce.Duration(r.electionTimeout, DefaultElectionTimeout))
}

// Start starts the raft node.
func (r *Raft) Start() error {
	r.infof("node starting")
	defer func() {
		r.infof("node started")
	}()

	r.server = NewServer().WithBindAddr(r.BindAddr()).WithLogger(r.log)
	r.server.SetAppendEntriesHandler(r.handleAppendEntries)
	r.server.SetRequestvoteHandler(r.handleRequestVote)

	r.infof("node rpc server starting, listening on: %s", r.BindAddr())
	err := r.server.Start()
	if err != nil {
		return err
	}

	r.infof("node beginning internal tickers")
	r.leaderCheckTick = worker.NewInterval(r.leaderCheck, DefaultLeaderCheckTick).WithDelay(r.RandomLeaderLeaseTimeout())
	r.heartbeatTick = worker.NewInterval(r.heartbeat, DefaultHeartbeatTick).WithDelay(r.RandomLeaderLeaseTimeout())
	r.leaderCheckTick.Start()
	r.heartbeatTick.Start()
	r.infof("leaderCheck start delay %v", r.leaderCheckTick.Delay())
	r.infof("heartbeatTick start delay %v", r.heartbeatTick.Delay())
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

// --------------------------------------------------------------------------------
// utility methods.
// --------------------------------------------------------------------------------

func (r *Raft) transitionTo(state FSMState) {
	defer func() { r.infof("transitioning to %s", state) }()

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

// LeaderCheck is the action that fires on a heartbeat to check if the leader lease
// has expired.
func (r *Raft) leaderCheck() error {
	if r.State() == FSMStateLeader {
		return nil
	}

	now := time.Now().UTC()
	if r.lastLeaderContact.IsZero() || now.Sub(r.lastLeaderContact) > r.RandomLeaderLeaseTimeout() {
		return r.election()
	}

	return nil
}

func (r *Raft) election() error {
	r.transitionTo(FSMStateCandidate)
	started := time.Now()
	for {
		// if we've been bumped out of candidate state,
		// stop the election cycle.
		if r.State() != FSMStateCandidate {
			return nil
		}

		if time.Since(started) > r.electionTimeout {
			r.votedFor = ""
			r.lastLeaderContact = time.Time{}
			r.transitionTo(FSMStateFollower)
		}

		if retry, err := r.requestVote(); err != nil {
			return err
		} else if !retry {
			return nil
		}
	}
}

func (r *Raft) requestVote() (retry bool, err error) {
	r.votedFor = ""
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

			res, err := c.RequestVote(&voteRequest)
			if err != nil {
				r.infof("requesting vote from %s: failed", c.RemoteAddr())
				errs <- err
			} else {
				r.infof("requesting vote from %s: %v", c.RemoteAddr(), res.Granted)
				results <- res
			}
		}(peer)
	}
	wg.Wait()

	if totalErrs := len(errs); totalErrs > 0 {
		for index := 0; index < totalErrs; index++ {
			r.err(<-errs)
		}
	}

	r.infof("election complete")
	result := r.countVotes(results)
	if result == 1 { // we're now the leader.
		r.votedFor = r.id
		r.transitionTo(FSMStateLeader)
		return
	} else if result == -1 {
		r.votedFor = ""
		r.lastLeaderContact = time.Time{}
		r.transitionTo(FSMStateFollower)
		return
	}

	// we tied, try again
	retry = true
	return
}

// countVotes returns the aggregate votes for in an election from rpc responses.
// it returns and integer, indicating victory, tie, or loss.
//  1 == victory
//  0 == tie
// -1 == loss
func (r *Raft) countVotes(results chan *RequestVoteResults) int {
	// tabulate results
	total := len(results)
	votesFor := 0

	defer func() {
		r.infof("election results: %d votes for, %d total", votesFor, total)
	}()

	if total == 0 {
		return 1
	}

	for index := 0; index < total; index++ {
		result := <-results

		if result.Granted {
			votesFor = votesFor + 1
		}
	}

	majority := total >> 1
	if total%2 == 0 {
		if votesFor > majority {
			return 1
		} else if votesFor == majority {
			return 0
		}
		return -1
	}

	if votesFor > majority {
		return 1
	}
	return -1
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
	return nil
}

func (r *Raft) handleAppendEntries(args *AppendEntries, res *AppendEntriesResults) error {
	if args.Term < r.currentTerm {
		*res = AppendEntriesResults{
			Success: false,
			Term:    r.currentTerm,
		}
		return nil
	}

	if r.lastLeaderContact.IsZero() {
		r.infof("first leader contact")
	}

	// advance the term ...
	r.currentTerm = args.Term
	r.lastLeaderContact = time.Now().UTC()
	*res = AppendEntriesResults{
		Success: true,
		Term:    r.currentTerm,
	}

	return nil
}

func (r *Raft) handleRequestVote(args *RequestVote, res *RequestVoteResults) error {
	if args.Term < r.currentTerm {
		r.infof("node is failing handle requestVote from %s because args.Term(%d) < currentTerm(%d)", args.Candidate, args.Term, r.currentTerm)
		*res = RequestVoteResults{
			Term:    r.currentTerm,
			Granted: false,
		}
		return nil
	}

	if len(r.votedFor) > 0 && r.votedFor != args.Candidate {
		*res = RequestVoteResults{
			Term:    args.Term,
			Granted: false,
		}
		return nil
	}

	// kill open election attempts.
	r.transitionTo(FSMStateFollower)

	r.lastVoteGranted = time.Now().UTC()
	r.votedFor = args.Candidate
	*res = RequestVoteResults{
		Term:    args.Term,
		Granted: true,
	}
	return nil
}

func (r *Raft) dialPeers() error {
	for _, peer := range r.peers {
		if err := peer.Open(); err != nil {
			return err
		}
	}
	return nil
}

func (r *Raft) infof(format string, args ...interface{}) {
	if r.log != nil {
		r.log.SubContext("raft").SubContext(r.id).SubContext(fmt.Sprintf("%v", r.State())).SyncInfof(format, args...)
	}
}

func (r *Raft) err(err error) {
	if r.log != nil {
		r.log.SubContext("raft").SubContext(r.id).SubContext(fmt.Sprintf("%v", r.State())).SyncTrigger(logger.Errorf(logger.Error, "%v", err))
	}
}
