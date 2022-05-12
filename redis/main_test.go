/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package redis_test

import (
	"context"

	"github.com/mediocregopher/radix/v4"
)

// MockRadixClient implements radix.Client for testing.
type MockRadixClient struct {
	radix.Client
	Ops chan radix.Action
}

// Do implements part of the radix client interface.
func (mrc *MockRadixClient) Do(ctx context.Context, action radix.Action) error {
	pushDone := make(chan struct{})
	go func() {
		defer close(pushDone)
		mrc.Ops <- action
	}()
	select {
	case <-ctx.Done():
		return context.Canceled
	case <-pushDone:
		return nil
	}
}
