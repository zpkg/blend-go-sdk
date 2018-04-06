package oauth

import (
	"testing"

	assert "github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/util"
)

func TestDeserializeJWTToken(t *testing.T) {
	assert := assert.New(t)

	testToken, err := SerializeJWT(util.Crypto.MustCreateKey(32), &JWTPayload{AUD: "client_id"})
	assert.Nil(err)

	jwt, err := DeserializeJWT(testToken)
	assert.Nil(err)
	assert.NotNil(jwt)
	assert.NotEmpty(jwt.Payload.AUD)
}
