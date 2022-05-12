/*

Copyright (c) 2022 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package crypto

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestEncryptDecrypt(t *testing.T) {
	assert := assert.New(t)
	key, err := CreateKey(32)
	assert.Nil(err)
	plaintext := "Mary Jane Hawkins"

	ciphertext, err := Encrypt(key, []byte(plaintext))
	assert.Nil(err)

	decryptedPlaintext, err := Decrypt(key, ciphertext)
	assert.Nil(err)
	assert.Equal(plaintext, string(decryptedPlaintext))
}
