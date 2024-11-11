/*

Copyright (c) 2022 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package envoyutil_test

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	sdkAssert "github.com/zpkg/blend-go-sdk/assert"
	"github.com/zpkg/blend-go-sdk/r2"
	"github.com/zpkg/blend-go-sdk/web"

	"github.com/zpkg/blend-go-sdk/envoyutil"
)

func TestWithClientIdentity(t *testing.T) {
	assert := sdkAssert.New(t)

	ctx := context.Background()
	newCtx := envoyutil.WithClientIdentity(ctx, "web.site")
	assert.Empty(envoyutil.GetClientIdentity(ctx))
	assert.Equal("web.site", envoyutil.GetClientIdentity(newCtx))
}

func TestClientIdentityRequired(t *testing.T) {
	assert := sdkAssert.New(t)

	app := web.MustNew()
	var capturedContext *web.Ctx
	cip := envoyutil.SPIFFEClientIdentityProvider(
		envoyutil.OptDeniedIdentities("gw.blend"),
	)
	verifier := envoyutil.SPIFFEServerIdentityProvider(
		envoyutil.OptAllowedIdentities("idea.blend"),
	)
	app.GET(
		"/",
		func(ctx *web.Ctx) web.Result {
			capturedContext = ctx
			return web.JSON.OK()
		},
		envoyutil.ClientIdentityRequired(cip, verifier),
		web.JSONProviderAsDefault,
	)

	body, meta, err := web.MockGet(app, "/").Bytes()
	assert.Nil(err)
	assert.Equal(http.StatusUnauthorized, meta.StatusCode, "Fail on missing header")
	assert.Nil(capturedContext)
	var expected error = &envoyutil.XFCCValidationError{Class: envoyutil.ErrMissingXFCC}
	invalidXFCCJSONEqual(assert, expected, body)

	xfcc := `""`
	body, meta, err = web.MockGet(app, "/", r2.OptHeaderValue(envoyutil.HeaderXFCC, xfcc)).Bytes()
	assert.Nil(err)
	assert.Equal(http.StatusBadRequest, meta.StatusCode, "Fail on empty header")
	assert.Nil(capturedContext)
	expected = &envoyutil.XFCCExtractionError{Class: envoyutil.ErrInvalidXFCC, XFCC: xfcc}
	invalidXFCCJSONEqual(assert, expected, body)

	xfcc = "something=bad"
	body, meta, err = web.MockGet(app, "/", r2.OptHeaderValue(envoyutil.HeaderXFCC, xfcc)).Bytes()
	assert.Nil(err)
	assert.Equal(http.StatusBadRequest, meta.StatusCode, "Fail on malformed header")
	assert.Nil(capturedContext)
	expected = &envoyutil.XFCCExtractionError{Class: envoyutil.ErrInvalidXFCC, XFCC: xfcc}
	invalidXFCCJSONEqual(assert, expected, body)

	xfcc = "By=spiffe://cluster.local/ns/blend/sa/idea;URI=spiffe://cluster.local/ns/blend/sa/should-end/sa/extra"
	body, meta, err = web.MockGet(app, "/", r2.OptHeaderValue(envoyutil.HeaderXFCC, xfcc)).Bytes()
	assert.Nil(err)
	assert.Equal(http.StatusBadRequest, meta.StatusCode, "Fail on unexpected SPIFFE format in `URI`")
	assert.Nil(capturedContext)
	expected = &envoyutil.XFCCExtractionError{Class: envoyutil.ErrInvalidClientIdentity, XFCC: xfcc}
	invalidXFCCJSONEqual(assert, expected, body)

	xfcc = `By=spiffe://cluster.local/ns/blend/sa/idea;Hash=468ed33be74eee6556d90c0149c1309e9ba61d6425303443c0748a02dd8de688;Subject="/C=US/ST=CA/L=San Francisco/OU=Lyft/CN=Test Client"`
	body, meta, err = web.MockGet(app, "/", r2.OptHeaderValue(envoyutil.HeaderXFCC, xfcc)).Bytes()
	assert.Nil(err)
	assert.Equal(http.StatusUnauthorized, meta.StatusCode, "Fail on missing client identity")
	assert.Nil(capturedContext)
	expected = &envoyutil.XFCCValidationError{Class: envoyutil.ErrInvalidClientIdentity, XFCC: xfcc}
	invalidXFCCJSONEqual(assert, expected, body)

	xfcc = `By=spiffe://cluster.local/ns/blend/sa/idea;Hash=468ed33be74eee6556d90c0149c1309e9ba61d6425303443c0748a02dd8de688;Subject="/C=US/ST=CA/L=San Francisco/OU=Lyft/CN=Test Client";URI=spiffe://cluster.local/ns/blend/sa/gw`
	body, meta, err = web.MockGet(app, "/", r2.OptHeaderValue(envoyutil.HeaderXFCC, xfcc)).Bytes()
	assert.Nil(err)
	assert.Equal(http.StatusUnauthorized, meta.StatusCode, "Fail on denied client identity")
	assert.Nil(capturedContext)
	expected = &envoyutil.XFCCValidationError{
		Class: envoyutil.ErrDeniedClientIdentity,
		XFCC:  xfcc,
		// NOTE: This should really be a `map[string]string`. We use a `map[string]interface{}`
		//       so that the comparison in `invalidXFCCJSONEqual()` passes.
		Metadata: map[string]interface{}{"clientIdentity": "gw.blend"},
	}
	invalidXFCCJSONEqual(assert, expected, body)

	xfcc = `By=spiffe://cluster.local/ns/blend/sa/idea;Hash=468ed33be74eee6556d90c0149c1309e9ba61d6425303443c0748a02dd8de688;Subject="/C=US/ST=CA/L=San Francisco/OU=Lyft/CN=Test Client";URI=spiffe://cluster.local/ns/blend/sa/twtr`
	meta, err = web.MockGet(app, "/", r2.OptHeaderValue(envoyutil.HeaderXFCC, xfcc)).Discard()
	assert.Nil(err)
	assert.Equal(http.StatusOK, meta.StatusCode, "Success on valid header")
	assert.NotNil(capturedContext)
	assert.Equal("twtr.blend", envoyutil.GetClientIdentity(capturedContext.Context()))
	capturedContext = nil

	xfcc = `By=mailto:John.Doe@example.com;URI=spiffe://cluster.local/ns/blend/sa/peas`
	body, meta, err = web.MockGet(app, "/", r2.OptHeaderValue(envoyutil.HeaderXFCC, xfcc)).Bytes()
	assert.Nil(err)
	assert.Equal(http.StatusBadRequest, meta.StatusCode, "Fail on invalid server identity")
	assert.Nil(capturedContext)
	expected = &envoyutil.XFCCExtractionError{
		Class: envoyutil.ErrInvalidServerIdentity,
		XFCC:  xfcc,
	}
	invalidXFCCJSONEqual(assert, expected, body)

	xfcc = `By=spiffe://cluster.local/ns/blend/sa/outside;URI=spiffe://cluster.local/ns/blend/sa/peas`
	body, meta, err = web.MockGet(app, "/", r2.OptHeaderValue(envoyutil.HeaderXFCC, xfcc)).Bytes()
	assert.Nil(err)
	assert.Equal(http.StatusUnauthorized, meta.StatusCode, "Fail on wrong server identity")
	assert.Nil(capturedContext)
	expected = &envoyutil.XFCCValidationError{
		Class: envoyutil.ErrDeniedServerIdentity,
		XFCC:  xfcc,
		// NOTE: This should really be a `map[string]string`. We use a `map[string]interface{}`
		//       so that the comparison in `invalidXFCCJSONEqual()` passes.
		Metadata: map[string]interface{}{"serverIdentity": "outside.blend"},
	}
	invalidXFCCJSONEqual(assert, expected, body)

	// Unrecoverable error: here we simulate `envoyutil` user error by using
	// `nil` for `cip`.
	app = web.MustNew()
	app.GET(
		"/",
		func(ctx *web.Ctx) web.Result {
			return web.JSON.OK()
		},
		envoyutil.ClientIdentityRequired(nil),
		web.JSONProviderAsDefault,
	)
	body, meta, err = web.MockGet(app, "/").Bytes()
	assert.Nil(err)
	assert.Equal(http.StatusInternalServerError, meta.StatusCode, "Fail on unrecoverable")
	assert.Equal("\"Internal Server Error\"\n", string(body))
}

func TestClientIdentityAware(t *testing.T) {
	assert := sdkAssert.New(t)

	app := web.MustNew()
	cip := envoyutil.SPIFFEClientIdentityProvider(
		envoyutil.OptDeniedIdentities("gw.blend"),
	)
	verifier := envoyutil.SPIFFEServerIdentityProvider(
		envoyutil.OptAllowedIdentities("quasar.blend"),
	)
	app.GET("/",
		func(_ *web.Ctx) web.Result {
			return web.JSON.OK()
		},
		envoyutil.ClientIdentityAware(cip, verifier),
		web.JSONProviderAsDefault,
	)

	meta, err := web.MockGet(app, "/").Discard()
	assert.Nil(err)
	assert.Equal(http.StatusOK, meta.StatusCode, "Don't fail on missing header")

	meta, err = web.MockGet(app, "/", r2.OptHeaderValue(envoyutil.HeaderXFCC, "something=bad")).Discard()
	assert.Nil(err)
	assert.Equal(http.StatusOK, meta.StatusCode, "Don't fail on malformed header")

	meta, err = web.MockGet(app, "/", r2.OptHeaderValue(envoyutil.HeaderXFCC, `By=spiffe://cluster.local/ns/blend/sa/quasar;Hash=468ed33be74eee6556d90c0149c1309e9ba61d6425303443c0748a02dd8de688;Subject="/C=US/ST=CA/L=San Francisco/OU=Lyft/CN=Test Client"`)).Discard()
	assert.Nil(err)
	assert.Equal(http.StatusOK, meta.StatusCode, "Don't fail on missing workload")

	meta, err = web.MockGet(app, "/", r2.OptHeaderValue(envoyutil.HeaderXFCC, `By=spiffe://cluster.local/ns/blend/sa/quasar;Hash=468ed33be74eee6556d90c0149c1309e9ba61d6425303443c0748a02dd8de688;Subject="/C=US/ST=CA/L=San Francisco/OU=Lyft/CN=Test Client";URI=spiffe://cluster.local/ns/blend/sa/gw`)).Discard()
	assert.Nil(err)
	assert.Equal(http.StatusOK, meta.StatusCode, "Don't fail on denied client identity")

	meta, err = web.MockGet(app, "/", r2.OptHeaderValue(envoyutil.HeaderXFCC, `By=spiffe://cluster.local/ns/blend/sa/quasar;Hash=468ed33be74eee6556d90c0149c1309e9ba61d6425303443c0748a02dd8de688;Subject="/C=US/ST=CA/L=San Francisco/OU=Lyft/CN=Test Client";URI=spiffe://cluster.local/ns/blend/sa/books`)).Discard()
	assert.Nil(err)
	assert.Equal(http.StatusOK, meta.StatusCode, "Success on valid header")

	// Unrecoverable error: here we simulate `envoyutil` user error by using
	// `nil` for `cip`.
	app = web.MustNew()
	app.GET(
		"/",
		func(ctx *web.Ctx) web.Result {
			return web.JSON.OK()
		},
		envoyutil.ClientIdentityAware(nil),
		web.JSONProviderAsDefault,
	)
	body, meta, err := web.MockGet(app, "/").Bytes()
	assert.Nil(err)
	assert.Equal(http.StatusInternalServerError, meta.StatusCode, "Fail on unrecoverable")
	assert.Equal("\"Internal Server Error\"\n", string(body))
}

func invalidXFCCJSONEqual(assert *sdkAssert.Assertions, expected error, actual []byte) {
	switch expected.(type) {
	case *envoyutil.XFCCExtractionError:
		unmarshaledActual := &envoyutil.XFCCExtractionError{}
		err := json.Unmarshal(actual, unmarshaledActual)
		assert.Nil(err)
		assert.Equal(expected, unmarshaledActual)
	case *envoyutil.XFCCValidationError:
		unmarshaledActual := &envoyutil.XFCCValidationError{}
		err := json.Unmarshal(actual, unmarshaledActual)
		assert.Nil(err)
		assert.Equal(expected, unmarshaledActual)
	default:
		assert.FailNow()
	}
}
