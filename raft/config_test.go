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
