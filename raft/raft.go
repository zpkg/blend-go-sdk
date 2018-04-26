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
		id:              uuid.V4().String(),
		state:           int32(FSMStateFollower),
		bindAddr:        DefaultBindAddr,
		leaderCheckTick: DefaultLeaderCheckTick,
		heartbeatTick:   DefaultHeartbeatTick,
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
	id       string
	log      *logger.Logger
	bindAddr string

	leaderLeaseTimeout time.Duration
	electionTimeout    time.Duration

	leaderCheckTick time.Duration
	heartbeatTick   time.Duration

	currentTerm       uint64
	currentLeader     string
	votedFor          string
	lastLeaderContact time.Time
	lastVoteGranted   time.Time

	// state is the current fsm state
	state int32

	server Server
	peers  []Client

	leaderCheckTicker   *worker.Interval
	sendHeartbeatTicker *worker.Interval

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

// Leader returns the current known leader.
func (r *Raft) Leader() string {
	if r.State() == FSMStateLeader {
		return r.id
	}
	return r.currentLeader
}

// Term returns the current raft term.
func (r *Raft) Term() uint64 {
	return r.currentTerm
}

// LastLeaderContact is the last time we heard from the leader.
func (r *Raft) LastLeaderContact() time.Time {
	return r.lastLeaderContact
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
func (r *Raft) WithPeer(peer Client) *Raft {
	r.peers = append(r.peers, peer)
	return r
}

// Peers returns the raft peers.
func (r *Raft) Peers() []Client {
	return r.peers
}

// WithServer sets the rpc server.
func (r *Raft) WithServer(server Server) *Raft {
	r.server = server
	return r
}

// Server returns the rpc server.
func (r *Raft) Server() Server {
	return r.server
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

// WithLeaderCheckTick sets the leader check tick.
func (r *Raft) WithLeaderCheckTick(d time.Duration) *Raft {
	r.leaderCheckTick = d
	return r
}

// LeaderCheckTick returns the leader check tick time.
func (r *Raft) LeaderCheckTick() time.Duration {
	return r.leaderCheckTick
}

// WithHeartbeatTick sets the heartbeat tick.
func (r *Raft) WithHeartbeatTick(d time.Duration) *Raft {
	r.heartbeatTick = d
	return r
}

// HeartbeatTick returns the heartbeat tick rate.
func (r *Raft) HeartbeatTick() time.Duration {
	return r.heartbeatTick
}

// Start starts the raft node.
func (r *Raft) Start() error {
	r.infof("node starting")
	defer func() {
		r.infof("node started")
	}()

	if len(r.peers) == 0 {
		r.infof("operating as single node configuration")
		r.transitionTo(FSMStateLeader)
		return nil
	}
	if r.server == nil {
		r.server = NewRPCServer().WithBindAddr(r.BindAddr()).WithLogger(r.log)
	}

	r.server.SetAppendEntriesHandler(r.handleAppendEntries)
	r.server.SetRequestVoteHandler(r.handleRequestVote)

	r.infof("node rpc server starting, listening on: %s", r.BindAddr())
	err := r.server.Start()
	if err != nil {
		return err
	}

	r.infof("node beginning internal tickers")
	r.leaderCheckTicker = worker.NewInterval(r.leaderCheck, r.leaderCheckTick).WithDelay(r.RandomLeaderLeaseTimeout())
	r.sendHeartbeatTicker = worker.NewInterval(r.sendHeartbeat, r.heartbeatTick).WithDelay(r.RandomLeaderLeaseTimeout())
	r.leaderCheckTicker.Start()
	r.sendHeartbeatTicker.Start()
	r.infof("leaderCheck start delay %v", r.leaderCheckTicker.Delay())
	r.infof("sendHeartbeatTick start delay %v", r.sendHeartbeatTicker.Delay())
	return nil
}

// Stop stops the node.
func (r *Raft) Stop() error {
	if r.leaderCheckTicker != nil {
		r.leaderCheckTicker.Stop()
		r.leaderCheckTicker = nil
	}
	if r.sendHeartbeatTicker != nil {
		r.sendHeartbeatTicker.Stop()
		r.sendHeartbeatTicker = nil
	}

	if r.server != nil {
		return r.server.Stop()
	}
	return nil
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
		go func(c Client) {
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
	result := r.processRequestVoteResults(results)
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

// processRequestVoteResults returns the aggregate votes for in an election from rpc responses.
func (r *Raft) processRequestVoteResults(results chan *RequestVoteResults) int {
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

	return r.voteOutcome(votesFor, total)
}

// voteOutcome compares votes for to total and  it returns and integer
// indicating victory, tie, or loss.
//  1 == victory
//  0 == tie
// -1 == loss
func (r *Raft) voteOutcome(votesFor, total int) int {
	if total == 0 {
		return 1
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

func (r *Raft) sendHeartbeat() error {
	if r.State() != FSMStateLeader {
		return nil
	}

	args := AppendEntries{
		Term:     r.currentTerm,
		LeaderID: r.id,
	}

	results := make(chan *AppendEntriesResults, len(r.peers))
	errs := make(chan error, len(r.peers))
	wg := sync.WaitGroup{}
	wg.Add(len(r.peers))

	for _, peer := range r.peers {
		go func(c Client) {
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

	if errCount := len(errs); errCount > 0 {
		for index := 0; index < errCount; index++ {
			r.err(<-errs)
		}
	}

	totalAnswers := len(results)
	successfulAnswers := 0
	latestTerm := r.currentTerm

	for index := 0; index < totalAnswers; index++ {
		answer := <-results
		if answer.Success {
			successfulAnswers = successfulAnswers + 1
		} else if answer.Term > latestTerm {
			latestTerm = answer.Term
		}
	}

	if r.voteOutcome(successfulAnswers, totalAnswers) == -1 {
		r.currentTerm = latestTerm
		r.votedFor = ""
		r.lastLeaderContact = time.Time{}
		r.transitionTo(FSMStateFollower)
	}

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

	r.currentTerm = args.Term
	r.currentLeader = args.LeaderID
	r.lastLeaderContact = time.Now().UTC()
	*res = AppendEntriesResults{
		Success: true,
		Term:    r.currentTerm,
	}

	return nil
}

func (r *Raft) handleRequestVote(args *RequestVote, res *RequestVoteResults) error {
	if args.Term < r.currentTerm {
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
		r.log.SubContext("raft").SubContext(fmt.Sprintf("%v", r.State())).Infof(format, args...)
	}
}

func (r *Raft) err(err error) {
	if r.log != nil {
		r.log.SubContext("raft").SubContext(fmt.Sprintf("%v", r.State())).Trigger(logger.Errorf(logger.Error, "%v", err))
	}
}
