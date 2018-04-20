package raft

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
