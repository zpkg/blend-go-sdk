/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package crypto

import (
	"bytes"
	"crypto/hmac"
	"crypto/rand"
	"strings"
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
)

func Test_LocalTransit_Encrypt_Decrypt(t *testing.T) {
	t.Parallel()
	its := assert.New(t)

	plaintext := "mary jane hawkins"

	m := NewLocalTransit(
		OptLocalTransitContextProvider(func() string {
			return time.Date(2019, 04, 15, 01, 02, 03, 04, time.UTC).Format("20060102")
		}),
		OptLocalTransitKey(MustCreateKey(32)),
	)
	prefix := m.ContextProvider()

	ciphertext := new(bytes.Buffer)
	its.Nil(m.Encrypt(ciphertext, bytes.NewReader([]byte(plaintext))))

	cipherBytes := ciphertext.Bytes()

	its.True(len(cipherBytes) > KeyVersionSize+IVSize+HashSize)
	its.True(strings.HasPrefix(string(cipherBytes), prefix+":"), "we should prefix ciphertext with the current date")

	output := new(bytes.Buffer)
	its.Nil(m.Decrypt(output, bytes.NewReader(ciphertext.Bytes())))
	its.Equal(plaintext, output.String())
}

func Test_Encrypt_Decrypt_large(t *testing.T) {
	t.Parallel()
	its := assert.New(t)

	m := NewLocalTransit(OptLocalTransitKey(MustCreateKey(32)))
	m.ContextProvider = func() string {
		return time.Date(2019, 04, 15, 01, 02, 03, 04, time.UTC).Format("20060102")
	}
	plaintext := make([]byte, 64*1024) // 64kb of data
	_, err := rand.Read(plaintext)
	its.Nil(err)

	ciphertext := new(bytes.Buffer)
	its.Nil(m.Encrypt(ciphertext, bytes.NewReader(plaintext)))

	output := new(bytes.Buffer)
	its.Nil(m.Decrypt(output, ciphertext))
	its.True(hmac.Equal(output.Bytes(), plaintext))
}
