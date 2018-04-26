package raft

import (
	"sync"
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/uuid"
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

func createTestNode() *Raft {
	return New().
		WithID(uuid.V4().String()).
		WithLeaderCheckTick(time.Microsecond).
		WithHeartbeatTick(time.Microsecond).
		WithElectionTimeout(time.Microsecond).
		WithLeaderLeaseTimeout(time.Microsecond).
		WithServer(NewMockServer())
}

func createCluster(nodeCount int) []*Raft {
	if nodeCount <= 0 {
		return nil
	}
	if nodeCount == 1 {
		return []*Raft{createTestNode()}
	}

	// create all the nodes
	peers := make([]*Raft, nodeCount)
	for index := 0; index < nodeCount; index++ {
		peers[index] = createTestNode()
	}

	// cross wire all the nodes
	for i := 0; i < nodeCount; i++ {
		for j := 0; j < nodeCount; j++ {
			if i != j {
				peers[i].WithPeer(NewMockTransport(peers[j].ID(), peers[j].Server()))
			}
		}
	}
	return peers
}

func TestRaftSingleNode(t *testing.T) {
	assert := assert.New(t)
	assert.StartTimeout(time.Millisecond)
	defer assert.EndTimeout()

	cluster := createCluster(1)[0]

	wg := sync.WaitGroup{}
	wg.Add(1)
	var didTransitionToLeader bool
	cluster.SetLeaderHandler(func() {
		defer wg.Done()
		didTransitionToLeader = true
	})
	go func() { cluster.Start() }()
	defer cluster.Stop()
	wg.Wait()
	assert.True(didTransitionToLeader)
}

func TestRaftCluster(t *testing.T) {
	assert := assert.New(t)
	assert.StartTimeout(time.Millisecond)
	defer assert.EndTimeout()

	cluster := createCluster(3)
	for _, node := range cluster {
		assert.Nil(node.Start())
		defer node.Stop()
	}

	// assert that the cluster stabilizes on leader.
}
