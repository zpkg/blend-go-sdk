/*

Copyright (c) 2022 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package crypto

import (
	"testing"

	"github.com/zpkg/blend-go-sdk/assert"
)

func Test_Encrypt_Decrypt(t *testing.T) {
	t.Parallel()

	its := assert.New(t)
	key, err := CreateKey(32)
	its.Nil(err)
	plaintext := "Mary Jane Hawkins"

	ciphertext, err := Encrypt(key, []byte(plaintext))
	its.Nil(err)

	decryptedPlaintext, err := Decrypt(key, ciphertext)
	its.Nil(err)
	its.Equal(plaintext, string(decryptedPlaintext))
}
