/*

Copyright (c) 2022 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package redis_test

import (
	"context"

	radix "github.com/mediocregopher/radix/v4"
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
