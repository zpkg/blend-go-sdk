package oauth

// Result is the final result of the oauth exchange.
// It is the user profile of the user and the state information.
type Result struct {
	Profile *Profile
	State   *State
}
