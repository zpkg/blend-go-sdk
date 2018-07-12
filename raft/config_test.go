package raft

import (
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
)

func TestConfig(t *testing.T) {
	assert := assert.New(t)

	assert.Equal("foo", (&Config{ID: "foo"}).GetID())
	assert.Equal(":1234", (&Config{BindAddr: ":1234"}).GetBindAddr())
	assert.Equal(DefaultBindAddr, (&Config{}).GetBindAddr())

	assert.Nil((&Config{}).GetPeers())
	assert.Equal([]string{"bar"}, (&Config{Peers: []string{"bar"}}).GetPeers())

	assert.Equal(time.Second, (&Config{ElectionTimeout: time.Second}).GetElectionTimeout())
}

func TestConfigGetPeersFiltered(t *testing.T) {
	assert := assert.New(t)

	unfiltered := (&Config{
		Peers: []string{
			"worker-0",
			"worker-1",
			"worker-2",
			"worker-3",
			"worker-4",
		},
	}).GetPeersFiltered()
	assert.Len(unfiltered, 5)

	filtered := (&Config{
		ExcludePeer: "worker-2",
		Peers: []string{
			"worker-0",
			"worker-1",
			"worker-2",
			"worker-3",
			"worker-4",
		},
	}).GetPeersFiltered()

	assert.Len(filtered, 4)
	assert.None(filtered, func(v interface{}) bool {
		return v.(string) == "worker-2"
	})
}
