/*

Copyright (c) 2022 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package async

import (
	"context"
	"sync"
	"testing"

	"github.com/zpkg/blend-go-sdk/assert"
)

func TestWorker(t *testing.T) {
	assert := assert.New(t)

	var didWork bool
	wg := sync.WaitGroup{}
	wg.Add(1)
	w := NewWorker(func(_ context.Context, obj interface{}) error {
		defer wg.Done()
		didWork = true
		assert.Equal("hello", obj)
		return nil
	})
	go func() { _ = w.Start() }()
	<-w.NotifyStarted()

	assert.True(w.Latch.IsStarted())
	w.Enqueue("hello")
	wg.Wait()
	assert.Nil(w.Stop())

	assert.False(w.Latch.IsStarted())
	assert.True(didWork)
}
