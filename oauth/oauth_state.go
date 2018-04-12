package oauth

import (
	"bytes"
	"encoding/gob"

	"github.com/blend/go-sdk/exception"
)

// DeserializeState deserializes the oauth state.
func DeserializeState(raw string) (*State, error) {
	corpus, err := Base64URLDecode(raw)
	if err != nil {
		return nil, exception.Wrap(err)
	}
	buffer := bytes.NewBuffer(corpus)
	var state State
	if err := gob.NewDecoder(buffer).Decode(&state); err != nil {
		return nil, exception.Wrap(err)
	}

	return &state, nil
}

// SerializeOAuthState serializes the oauth state.
func SerializeOAuthState(state *State) (output string, err error) {
	buffer := bytes.NewBuffer(nil)
	err = gob.NewEncoder(buffer).Encode(state)
	if err != nil {
		err = exception.Wrap(err)
		return
	}
	output = Base64URLEncode(buffer.Bytes())
	return
}

// State is the oauth state.
type State struct {
	Token       string
	Secure      string
	RedirectURL string
}
