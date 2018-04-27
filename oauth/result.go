package oauth

import "time"

// Result is the final result of the oauth exchange.
// It is the user profile of the user and the state information.
type Result struct {
	AccessToken  string
	TokenType    string
	RefreshToken string
	Expiry       time.Time

	Profile *Profile
	State   *State
}
