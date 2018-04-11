package oauth

import (
	"testing"

	assert "github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/env"
	"github.com/blend/go-sdk/util"
)

func TestNewConfigFromEnv(t *testing.T) {
	assert := assert.New(t)
	defer env.Restore()

	env.Env().Set("GOOGLE_SKIP_DOMAIN_VALIDATION", "false")
	env.Env().Set("GOOGLE_REDIRECT_URI", "https://app.com/oauth/google")
	env.Env().Set("GOOGLE_VALID_DOMAINS", "foo.com,bar.com")
	env.Env().Set("GOOGLE_CLIENT_ID", "foo")
	env.Env().Set("GOOGLE_CLIENT_SECRET", "bar")

	cfg := NewConfigFromEnv()
	assert.False(cfg.GetSkipDomainValidation())
	assert.Equal("foo", cfg.GetClientID())
	assert.Equal("bar", cfg.GetClientSecret())
	assert.Equal("https://app.com/oauth/google", cfg.GetRedirectURI())
	assert.Len(cfg.GetValidDomains(), 2)
	assert.Equal("foo.com", cfg.GetValidDomains()[0])
	assert.Equal("bar.com", cfg.GetValidDomains()[1])
}

func TestConfig(t *testing.T) {
	assert := assert.New(t)

	assert.True(Config{}.IsZero())
	assert.True(Config{ClientID: "foo"}.IsZero())
	assert.False(Config{ClientID: "foo", ClientSecret: "bar"}.IsZero())
}

func TestConfigGetSecret(t *testing.T) {
	assert := assert.New(t)

	secret, err := Config{}.GetSecret()
	assert.Nil(err)
	assert.Empty(secret, "zero config should return no secret")

	withSecret := &Config{Secret: Base64Encode(util.Crypto.MustCreateKey(32))}
	secret, err = withSecret.GetSecret()
	assert.Nil(err)
	assert.NotEmpty(secret, "secret should be base64 enoded")

	malformed := &Config{Secret: "|||||"}
	secret, err = malformed.GetSecret()
	assert.NotNil(err)
	assert.Empty(secret)
}
