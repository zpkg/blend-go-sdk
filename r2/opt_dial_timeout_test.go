/*

Copyright (c) 2022 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package r2

import (
	"net/http"
	"testing"
	"time"

	"github.com/zpkg/blend-go-sdk/assert"
)

func TestOptDialTimeout(t *testing.T) {
	assert := assert.New(t)

	r := New(TestURL, OptDial(OptDialTimeout(time.Second)))
	assert.NotNil(r.Client.Transport.(*http.Transport).DialContext)
}
