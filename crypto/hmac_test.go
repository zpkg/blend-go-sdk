/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package crypto

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestHMAC(t *testing.T) {
	assert := assert.New(t)
	key, err := CreateKey(128)
	assert.Nil(err)
	plaintext := "123-12-1234"
	assert.Equal(
		HMAC512(key, []byte(plaintext)),
		HMAC512(key, []byte(plaintext)),
	)
}
