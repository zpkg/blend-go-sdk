/*

Copyright (c) 2022 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package web

import (
	"context"
	"testing"
	"time"

	"github.com/zpkg/blend-go-sdk/assert"
)

func TestContextApp(t *testing.T) {
	assert := assert.New(t)

	app := MustNew()

	ctx := WithApp(context.Background(), app)
	assert.NotNil(GetApp(ctx))
	assert.Nil(GetApp(context.Background()))
}

func TestContextRequestStart(t *testing.T) {
	assert := assert.New(t)

	ts := time.Date(2020, 06, 02, 12, 11, 10, 9, time.UTC)
	ctx := WithRequestStarted(context.Background(), ts)
	assert.Equal(ts, GetRequestStarted(ctx))
}
