package oauth

import (
	"testing"

	assert "github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/util"
)

func TestGenerateSecret(t *testing.T) {
	assert := assert.New(t)

	key := GenerateSecret()
	decoded, err := Base64Decode(key)
	assert.Nil(err)

	cipherText, err := util.Crypto.Encrypt(decoded, []byte("foo"))
	assert.Nil(err)
	plaintext, err := util.Crypto.Decrypt(decoded, cipherText)
	assert.Nil(err)

	assert.Equal("foo", string(plaintext))
}

func TestBase64Encode(t *testing.T) {
	assert := assert.New(t)

	decoded, err := Base64Decode(Base64Encode([]byte("foo")))
	assert.Nil(err)
	assert.Equal("foo", string(decoded))
}

func TestMustBase64Decode(t *testing.T) {
	assert := assert.New(t)
	decoded := MustBase64Decode(Base64Encode([]byte("foo")))
	assert.Equal("foo", string(decoded))
}
