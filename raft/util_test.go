package raft

import (
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
)

func TestRandomTimeout(t *testing.T) {
	assert := assert.New(t)

	v := 5 * time.Second
	assert.True(RandomTimeout(v) > v)
	assert.True(RandomTimeout(v) < 3*v)
}
