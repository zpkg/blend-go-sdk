/*

Copyright (c) 2022 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package r2

import (
	"net/http"
	"testing"

	"github.com/zpkg/blend-go-sdk/assert"
)

func TestOptTLSClientCert(t *testing.T) {
	assert := assert.New(t)

	r := New("https://foo.com", OptTLSClientCert(clientCert, clientKey))
	assert.NotNil(r.Client)
	assert.NotNil(r.Client.Transport)
	assert.NotNil(r.Client.Transport.(*http.Transport).TLSClientConfig)
	assert.NotEmpty(r.Client.Transport.(*http.Transport).TLSClientConfig.Certificates)
}

func TestOptTLSClientCertErrors(t *testing.T) {
	assert := assert.New(t)

	r := New("https://foo.com", OptTLSClientCert(nil, nil))
	assert.NotNil(r.Err)
}
