package raft

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestRaftCountVotes(t *testing.T) {
	assert := assert.New(t)

	r := New()

	empty := make(chan *RequestVoteResults, 0)
	assert.Equal(1, r.countVotes(empty))

	loss := make(chan *RequestVoteResults, 1)
	loss <- &RequestVoteResults{Granted: false}
	assert.Equal(-1, r.countVotes(loss))

	win := make(chan *RequestVoteResults, 1)
	win <- &RequestVoteResults{Granted: true}
	assert.Equal(1, r.countVotes(win))

	tie := make(chan *RequestVoteResults, 2)
	tie <- &RequestVoteResults{Granted: true}
	tie <- &RequestVoteResults{Granted: false}
	assert.Equal(0, r.countVotes(tie))

	twoLoss := make(chan *RequestVoteResults, 2)
	twoLoss <- &RequestVoteResults{Granted: false}
	twoLoss <- &RequestVoteResults{Granted: false}
	assert.Equal(-1, r.countVotes(twoLoss))

	twoWin := make(chan *RequestVoteResults, 2)
	twoWin <- &RequestVoteResults{Granted: true}
	twoWin <- &RequestVoteResults{Granted: true}
	assert.Equal(1, r.countVotes(twoWin))

	threeLoss := make(chan *RequestVoteResults, 3)
	threeLoss <- &RequestVoteResults{Granted: false}
	threeLoss <- &RequestVoteResults{Granted: false}
	threeLoss <- &RequestVoteResults{Granted: false}
	assert.Equal(-1, r.countVotes(threeLoss))

	threeOneWin := make(chan *RequestVoteResults, 3)
	threeOneWin <- &RequestVoteResults{Granted: true}
	threeOneWin <- &RequestVoteResults{Granted: false}
	threeOneWin <- &RequestVoteResults{Granted: false}
	assert.Equal(-1, r.countVotes(threeOneWin))

	threeTwoWin := make(chan *RequestVoteResults, 3)
	threeTwoWin <- &RequestVoteResults{Granted: true}
	threeTwoWin <- &RequestVoteResults{Granted: true}
	threeTwoWin <- &RequestVoteResults{Granted: false}
	assert.Equal(1, r.countVotes(threeTwoWin))

	threeThreeWin := make(chan *RequestVoteResults, 3)
	threeThreeWin <- &RequestVoteResults{Granted: true}
	threeThreeWin <- &RequestVoteResults{Granted: true}
	threeThreeWin <- &RequestVoteResults{Granted: true}
	assert.Equal(1, r.countVotes(threeTwoWin))

	fourLoss := make(chan *RequestVoteResults, 4)
	fourLoss <- &RequestVoteResults{Granted: false}
	fourLoss <- &RequestVoteResults{Granted: false}
	fourLoss <- &RequestVoteResults{Granted: false}
	fourLoss <- &RequestVoteResults{Granted: false}
	assert.Equal(1, r.countVotes(threeTwoWin))
}
