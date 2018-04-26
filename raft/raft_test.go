package raft

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestRaftCountVotes(t *testing.T) {
	assert := assert.New(t)

	r := New()

	assert.Equal(1, r.voteOutcome(0, 0))
	assert.Equal(1, r.voteOutcome(1, 0))

	assert.Equal(-1, r.voteOutcome(0, 1))
	assert.Equal(1, r.voteOutcome(1, 1))

	assert.Equal(-1, r.voteOutcome(0, 2))
	assert.Equal(0, r.voteOutcome(1, 2))
	assert.Equal(1, r.voteOutcome(2, 2))

	assert.Equal(-1, r.voteOutcome(0, 3))
	assert.Equal(-1, r.voteOutcome(1, 3))
	assert.Equal(1, r.voteOutcome(2, 3))
	assert.Equal(1, r.voteOutcome(3, 3))

	assert.Equal(-1, r.voteOutcome(0, 4))
	assert.Equal(-1, r.voteOutcome(1, 4))
	assert.Equal(0, r.voteOutcome(2, 4))
	assert.Equal(1, r.voteOutcome(3, 4))
	assert.Equal(1, r.voteOutcome(4, 4))

	assert.Equal(-1, r.voteOutcome(0, 4))
	assert.Equal(-1, r.voteOutcome(1, 4))
	assert.Equal(0, r.voteOutcome(2, 4))
	assert.Equal(1, r.voteOutcome(3, 4))
	assert.Equal(1, r.voteOutcome(4, 4))
}

func createSingleNode() *Raft {
	return New().WithServer(NewMockTransport
}

func TestRaftSingleNode(t *testing.T) {
	assert := assert.New(t)

}
