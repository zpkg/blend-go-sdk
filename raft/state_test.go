package raft

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestFSMState(t *testing.T) {
	assert := assert.New(t)

	assert.Equal("unset", Unset.String())
	assert.Equal("follower", Follower.String())
	assert.Equal("candidate", Candidate.String())
	assert.Equal("leader", Leader.String())
	assert.Equal("unknown", State(123).String())
}
