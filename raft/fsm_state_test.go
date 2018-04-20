package raft

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestFSMState(t *testing.T) {
	assert := assert.New(t)

	assert.Equal("unset", FSMStateUnset.String())
	assert.Equal("follower", FSMStateFollower.String())
	assert.Equal("candidate", FSMStateCandidate.String())
	assert.Equal("leader", FSMStateLeader.String())
	assert.Equal("unknown", FSMState(123).String())
}
