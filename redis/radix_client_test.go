/*

Copyright (c) 2022 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package redis_test

import (
	"bytes"
	"context"
	"testing"
	"time"

	radix "github.com/mediocregopher/radix/v4"

	"github.com/zpkg/blend-go-sdk/assert"
	"github.com/zpkg/blend-go-sdk/logger"
	"github.com/zpkg/blend-go-sdk/redis"
)

func Test_RadixClient_Do(t *testing.T) {
	its := assert.New(t)

	buf := new(bytes.Buffer)
	log := logger.Memory(buf)
	defer log.Close()

	logEvents := make(chan redis.Event)
	log.Listen("test", "test", redis.NewEventListener(func(_ context.Context, e redis.Event) {
		logEvents <- e
	}))

	mockRadixClient := &MockRadixClient{
		Ops: make(chan radix.Action, 1),
	}

	rc := &redis.RadixClient{
		Log:    log,
		Client: mockRadixClient,
	}

	var foo string
	its.Nil(rc.Do(context.TODO(), &foo, "GET", "foo"))
}

func Test_RadixClient_Do_timeout(t *testing.T) {
	its := assert.New(t)

	mockRadixClient := &MockRadixClient{
		Ops: make(chan radix.Action),
	}
	rc := &redis.RadixClient{
		Config: redis.Config{
			Timeout: time.Millisecond,
		},
		Client: mockRadixClient,
	}
	var foo string
	its.NotNil(rc.Do(context.Background(), &foo, "GET", "foo"))
}
