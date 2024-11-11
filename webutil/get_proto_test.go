/*

Copyright (c) 2022 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package webutil

import (
	"net/http"
	"testing"

	"github.com/zpkg/blend-go-sdk/assert"
)

func TestGetProto(t *testing.T) {
	assert := assert.New(t)

	headers := http.Header{}
	headers.Set(HeaderXForwardedProto, SchemeHTTPS)
	r := http.Request{
		Proto:  SchemeHTTP + "/1.0",
		Header: headers,
	}
	assert.Equal(SchemeHTTPS, GetProto(&r))

	headers = http.Header{}
	headers.Set(HeaderXForwardedProto, SchemeSPDY+","+SchemeHTTPS)
	r = http.Request{
		Proto:  SchemeHTTP + "/1.0",
		Header: headers,
	}
	assert.Equal(SchemeHTTPS, GetProto(&r))

	headers = http.Header{}
	r = http.Request{
		Header: headers,
	}
	assert.Empty(GetProto(&r))
}
