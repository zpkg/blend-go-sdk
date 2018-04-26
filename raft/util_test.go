package raft

import (
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
)

func TestRandomTimeout(t *testing.T) {
	assert := assert.New(t)

	v := 5 * time.Second
	assert.True(randomTimeout(v) > v)
	assert.True(randomTimeout(v) < 3*v)
}
