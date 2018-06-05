package exception

import "encoding/json"

// Error is a string wrapper that implements `error`.
// Use this to implement constant exception causes.
type Error string

// Error implements `error`.
func (e Error) Error() string {
	return string(e)
}

// MarshalJSON implements json.Marshaler.
func (e Error) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(e))
}
