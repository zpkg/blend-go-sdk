package cron

import (
	"context"
	"sync"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestIsJobCanceled(t *testing.T) {
	assert := assert.New(t)

	ctx, cancel := context.WithCancel(context.Background())
	assert.False(IsJobCancelled(ctx))

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		assert.True(IsJobCancelled(ctx))
	}()
	cancel()
	wg.Wait()
}
