/*

Copyright (c) 2022 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package certutil

import (
	"os"
	"testing"

	"github.com/zpkg/blend-go-sdk/assert"
)

func Test_JoinPEMs(t *testing.T) {
	its := assert.New(t)

	ca, err := os.ReadFile("testdata/ca.cert.pem")
	its.Nil(err)

	serverPartial, err := os.ReadFile("testdata/server.partial.cert.pem")
	its.Nil(err)

	serverFull, err := os.ReadFile("testdata/server.cert.pem")
	its.Nil(err)

	serverJoined := JoinPEMs(string(serverPartial), " ", string(ca))

	its.Equal(string(serverFull), serverJoined)
}
