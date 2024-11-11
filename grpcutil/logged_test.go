/*

Copyright (c) 2022 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package grpcutil

import (
	"context"
	"sync"
	"testing"

	"google.golang.org/grpc"

	"github.com/zpkg/blend-go-sdk/assert"
	"github.com/zpkg/blend-go-sdk/logger"
)

func TestLoggedClientUnary(t *testing.T) {
	assert := assert.New(t)

	log := logger.All()
	defer log.Close()

	events := make(chan RPCEvent, 1)
	var wg sync.WaitGroup
	wg.Add(1)
	log.Listen(FlagRPC, "test", func(_ context.Context, e logger.Event) {
		wg.Done()
		events <- e.(RPCEvent)
	})
	interceptor := LoggedClientUnary(log)

	//func(ctx context.Context, method string, req interface{}, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	err := interceptor(context.TODO(), "/example-string/v1/dog", "treats", nil, nil, grpc.UnaryInvoker(func(_ context.Context, _ string, _, _ interface{}, _ *grpc.ClientConn, opts ...grpc.CallOption) error {
		return nil
	}))
	assert.Nil(err)

	wg.Wait()

	assert.NotEmpty(events)
	got := <-events
	assert.Equal("/example-string/v1/dog", got.Method)
}
