package cron

import (
	"context"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestIsJobCanceled(t *testing.T) {
	assert := assert.New(t)

	ctx, cancel := context.WithCancel(context.Background())
	assert.False(IsJobCancelled(ctx))
	cancel()
	assert.True(IsJobCancelled(ctx))
}
