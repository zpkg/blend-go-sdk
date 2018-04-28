package raft

// ElectionOutcome is an election outcome.
type ElectionOutcome int

// String returns the string value for the outcome.
func (eo ElectionOutcome) String() string {
	switch eo {
	case ElectionVictory:
		return "victory"
	case ElectionTie:
		return "tie"
	case ElectionLoss:
		return "loss"
	default:
		return "unknown"
	}
}

const (
	// ElectionVictory is an election outcome.
	ElectionVictory ElectionOutcome = 1
	// ElectionTie is an election outcome.
	ElectionTie ElectionOutcome = 0
	// ElectionLoss is an election outcome.
	ElectionLoss ElectionOutcome = -1
)
