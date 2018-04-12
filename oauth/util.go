package oauth

import (
	"encoding/base64"

	"github.com/blend/go-sdk/exception"
	"github.com/blend/go-sdk/util"
)

// GenerateSecret is a helper to create secret keys.
func GenerateSecret() string {
	return Base64Encode(util.Crypto.MustCreateKey(32))
}

// Base64Decode decodes a string as base64.
func Base64Decode(corpus string) ([]byte, error) {
	contents, err := base64.StdEncoding.DecodeString(corpus)
	if err != nil {
		return nil, exception.Wrap(err)
	}
	return contents, nil
}

// Base64URLDecode decodes a string as base64.
func Base64URLDecode(corpus string) ([]byte, error) {
	contents, err := base64.URLEncoding.DecodeString(corpus)
	if err != nil {
		return nil, exception.Wrap(err)
	}
	return contents, nil
}

// MustBase64Decode decodes a string as base64 and panics if there is an error.
func MustBase64Decode(corpus string) []byte {
	contents, err := Base64Decode(corpus)
	if err != nil {
		panic(err)
	}
	return contents
}

// MustBase64URLDecode decodes a string as base64 and panics if there is an error.
func MustBase64URLDecode(corpus string) []byte {
	contents, err := Base64URLDecode(corpus)
	if err != nil {
		panic(err)
	}
	return contents
}

// Base64Encode encodes binary as a base64 string.
func Base64Encode(corpus []byte) string {
	return base64.StdEncoding.EncodeToString(corpus)
}

// Base64URLEncode encodes binary as a base64 string.
func Base64URLEncode(corpus []byte) string {
	return base64.URLEncoding.EncodeToString(corpus)
}
