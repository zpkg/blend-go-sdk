// Package raft is a very spare implementation of the raft consensus protocol concentrating on leader elections.
// It uses HTTP for rpc, and is designed to be used in an unreliable network environment.
// It is *extremely* experimental at the moment, and not designed for production use.
package raft
