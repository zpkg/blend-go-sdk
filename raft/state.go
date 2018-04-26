package raft

// State is a raft fsm state.
type State int

func (fsm State) String() string {
	switch fsm {
	case Unset:
		return "unset"
	case Follower:
		return "follower"
	case Candidate:
		return "candidate"
	case Leader:
		return "leader"
	default:
		return "unknown"
	}
}

const (
	// Unset is the unset fsm state.
	Unset State = 0
	// Follower is the follower state.
	Follower State = 1
	// Candidate is the follower state.
	Candidate State = 2
	// Leader is the follower state.
	Leader State = 3
)
