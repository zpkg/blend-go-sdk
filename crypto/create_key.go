/*

Copyright (c) 2022 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package crypto

import (
	cryptorand "crypto/rand"
	"encoding/base64"
	"encoding/hex"

	"github.com/zpkg/blend-go-sdk/ex"
)

// MustCreateKey creates a key, if an error is returned, it panics.
func MustCreateKey(keySize int) []byte {
	key, err := CreateKey(keySize)
	if err != nil {
		panic(err)
	}
	return key
}

// CreateKey creates a key of a given size by reading that much data off the crypto/rand reader.
func CreateKey(keySize int) ([]byte, error) {
	key := make([]byte, keySize)
	_, err := cryptorand.Read(key)
	if err != nil {
		return nil, err
	}
	return key, nil
}

// MustCreateKeyString generates a new key and returns it as a hex string.
func MustCreateKeyString(keySize int) string {
	return hex.EncodeToString(MustCreateKey(keySize))
}

// MustCreateKeyBase64String generates a new key and returns it as a base64 std encoding string.
func MustCreateKeyBase64String(keySize int) string {
	return base64.StdEncoding.EncodeToString(MustCreateKey(keySize))
}

// CreateKeyString generates a new key and returns it as a hex string.
func CreateKeyString(keySize int) (string, error) {
	key, err := CreateKey(keySize)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(key), nil
}

// CreateKeyBase64String generates a new key and returns it as a base64 std encoding string.
func CreateKeyBase64String(keySize int) (string, error) {
	key, err := CreateKey(keySize)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(key), nil
}

// ParseKey parses a key from a string.
func ParseKey(key string) ([]byte, error) {
	decoded, err := hex.DecodeString(key)
	if err != nil {
		return nil, ex.New(err)
	}
	if len(decoded) != DefaultKeySize {
		return nil, ex.New("parse key; invalid key length")
	}
	return decoded, nil
}
