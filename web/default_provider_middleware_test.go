/*

Copyright (c) 2022 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package web

import (
	"bytes"
	"testing"

	"github.com/zpkg/blend-go-sdk/assert"
	"github.com/zpkg/blend-go-sdk/webutil"
)

func TestDefaultProviderMiddlewares(t *testing.T) {
	assert := assert.New(t)

	r := applyMiddleware(JSONProviderAsDefault)
	_, ok := r.DefaultProvider.(JSONResultProvider)
	assert.True(ok)

	r = applyMiddleware(ViewProviderAsDefault)
	_, ok = r.DefaultProvider.(*ViewCache)
	assert.True(ok)

	r = applyMiddleware(XMLProviderAsDefault)
	_, ok = r.DefaultProvider.(XMLResultProvider)
	assert.True(ok)

	r = applyMiddleware(TextProviderAsDefault)
	_, ok = r.DefaultProvider.(TextResultProvider)
	assert.True(ok)
}

func applyMiddleware(middleware Middleware) (output *Ctx) {
	middleware(func(ctx *Ctx) Result {
		output = ctx
		return NoContent
	})(NewCtx(webutil.NewMockResponse(new(bytes.Buffer)), webutil.NewMockRequest("GET", "/")))
	return
}
