package oauth

import (
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/request"
	"github.com/blend/go-sdk/util"
	"github.com/blend/go-sdk/uuid"
)

func TestNewFromConfig(t *testing.T) {
	assert := assert.New(t)

	m, err := NewFromConfig(&Config{
		SkipDomainValidation: false,
		RedirectURI:          "https://app.com/oauth/google",
		ValidDomains:         []string{"foo.com", "bar.com"},
		ClientID:             "foo_client",
		ClientSecret:         "bar_secret",
	})

	assert.Nil(err)
	assert.NotEmpty(m.Secret())
	assert.Equal("https://app.com/oauth/google", m.RedirectURI())
	assert.Len(m.ValidDomains(), 2)
	assert.Equal("foo_client", m.ClientID())
	assert.Equal("bar_secret", m.ClientSecret())
}

func TestNewFromConfigWithSecret(t *testing.T) {
	assert := assert.New(t)

	m, err := NewFromConfig(&Config{
		Secret: Base64Encode([]byte("test string")),
	})

	assert.Nil(err)
	assert.NotEmpty(m.Secret())
	assert.Equal("test string", string(m.Secret()))
}

func TestManagerValidDomains(t *testing.T) {
	assert := assert.New(t)

	assert.Empty(New().ValidDomains())

	domains := New().WithHostedDomain("foo.com").ValidDomains()
	assert.NotEmpty(domains)
	assert.Len(domains, 1)
	assert.Equal("foo.com", domains[0])

	domains = New().WithValidDomains("foo.com", "bar.com").ValidDomains()
	assert.NotEmpty(domains)
	assert.Len(domains, 2)
	assert.Any(domains, func(v interface{}) bool { return v.(string) == "foo.com" })
	assert.Any(domains, func(v interface{}) bool { return v.(string) == "bar.com" })

	domains = New().WithHostedDomain("buzz.com").WithValidDomains("foo.com", "bar.com").ValidDomains()
	assert.NotEmpty(domains)
	assert.Len(domains, 3)
	assert.Any(domains, func(v interface{}) bool { return v.(string) == "foo.com" })
	assert.Any(domains, func(v interface{}) bool { return v.(string) == "bar.com" })
	assert.Any(domains, func(v interface{}) bool { return v.(string) == "buzz.com" })

	domains = New().WithHostedDomain("bar.com").WithValidDomains("foo.com", "bar.com").ValidDomains()
	assert.NotEmpty(domains)
	assert.Len(domains, 2)
	assert.Any(domains, func(v interface{}) bool { return v.(string) == "foo.com" })
	assert.Any(domains, func(v interface{}) bool { return v.(string) == "bar.com" })
}

func TestManagerOAuthURL(t *testing.T) {
	assert := assert.New(t)

	m := New().
		WithClientID("test_client_id").
		WithHostedDomain("test.blend.com").
		WithRedirectURI("https://local.shortcut-service.centrio.com/oauth/google")

	_, err := m.OAuthURL()
	assert.Nil(err)
}

func TestManagerOAuthURLSecurity(t *testing.T) {
	assert := assert.New(t)

	m := New().
		WithSecret(util.Crypto.MustCreateKey(32)).
		WithClientID("test_client_id").
		WithRedirectURI("https://local.shortcut-service.centrio.com/oauth/google")

	urlFragment, err := m.OAuthURL()
	assert.Nil(err)

	u, err := url.Parse(urlFragment)
	assert.Nil(err)
	assert.NotEmpty(u.Query().Get("state"))
	assert.Equal("test_client_id", u.Query().Get("client_id"))

	state := u.Query().Get("state")
	deserialized, err := DeserializeState(state)
	assert.Nil(err)
	assert.Nil(m.ValidateOAuthState(deserialized))
}

func TestManagerOAuthURLRedirect(t *testing.T) {
	assert := assert.New(t)

	m := New().
		WithSecret(util.Crypto.MustCreateKey(32)).
		WithClientID("test_client_id").
		WithRedirectURI("https://local.shortcut-service.centrio.com/oauth/google")

	urlFragment, err := m.OAuthURL("bar_foo")
	assert.Nil(err)

	u, err := url.Parse(urlFragment)
	assert.Nil(err)
	assert.NotEmpty(u.Query().Get("state"))

	state := u.Query().Get("state")
	deserialized, err := DeserializeState(state)
	assert.Nil(err)
	assert.Nil(m.ValidateOAuthState(deserialized))
	assert.Equal("bar_foo", deserialized.RedirectURL)
}

func TestManagerValidateProfile(t *testing.T) {
	assert := assert.New(t)

	blender := &Profile{
		Email: "bailey@blend.com",
	}

	personal := &Profile{
		Email: "bailey@gmail.com",
	}

	suffixMatch := &Profile{
		Email: "bailey@sailblend.com",
	}

	prefixMatch := &Profile{
		Email: "bailey@blend.com.au",
	}

	empty := New()
	assert.NotNil(empty.ValidateProfile(blender), "we should error if no domains are configured with defaults")

	unvalidated := New().WithSkipDomainValidation(true)
	assert.Nil(unvalidated.ValidateProfile(blender), "we should not error if skip validation is true")
	assert.Nil(unvalidated.ValidateProfile(personal), "we should not error if skip validation is true")
	assert.Nil(unvalidated.ValidateProfile(suffixMatch), "we should not error if skip validation is true")
	assert.Nil(unvalidated.ValidateProfile(prefixMatch), "we should not error if skip validation is true")

	validated := New().WithValidDomains("blend.com")
	assert.Nil(validated.ValidateProfile(blender), "we should pass for @blend.com")
	assert.NotNil(validated.ValidateProfile(personal), "we fail for non-@blend.com emails")
	assert.NotNil(validated.ValidateProfile(suffixMatch), "we fail for non-@blend.com emails")
	assert.NotNil(validated.ValidateProfile(prefixMatch), "we fail for non-@blend.com emails")

	hostedOnly := New().WithHostedDomain("blend.com")
	assert.Nil(hostedOnly.ValidateProfile(blender), "we should pass for @blend.com")
	assert.NotNil(hostedOnly.ValidateProfile(personal), "we fail for non-@blend.com emails")
	assert.NotNil(hostedOnly.ValidateProfile(suffixMatch), "we fail for non-@blend.com emails")
	assert.NotNil(hostedOnly.ValidateProfile(prefixMatch), "we fail for non-@blend.com emails")

	validatedPrefixed := New().WithValidDomains("@blend.com")
	assert.Nil(validatedPrefixed.ValidateProfile(blender), "we should pass for @blend.com")
	assert.NotNil(validatedPrefixed.ValidateProfile(personal), "we fail for non-@blend.com emails")
	assert.NotNil(validatedPrefixed.ValidateProfile(suffixMatch), "we fail for non-@blend.com emails")
	assert.NotNil(validatedPrefixed.ValidateProfile(prefixMatch), "we fail for non-@blend.com emails")

	multi := New().WithValidDomains("@blend.com", "sailblend.com")
	assert.Nil(multi.ValidateProfile(blender), "we should pass for @blend.com")
	assert.NotNil(multi.ValidateProfile(personal), "we fail for non-@blend.com emails")
	assert.Nil(multi.ValidateProfile(suffixMatch), "we should pass for @sailblend.com as it was configured too")
	assert.NotNil(multi.ValidateProfile(prefixMatch), "we fail for non-@blend.com emails")
}

func TestManagerValidateOAuthState(t *testing.T) {
	assert := assert.New(t)

	m := New().WithSecret(util.Crypto.MustCreateKey(32))
	assert.Nil(m.ValidateOAuthState(m.CreateState()))
}

func TestManagerState(t *testing.T) {
	assert := assert.New(t)

	m := New().WithSecret(util.Crypto.MustCreateKey(32))
	serialized, err := m.State()
	assert.Nil(err)
	assert.NotEmpty(serialized, "at the baseline we should have the anti-forgery pair")

	state := m.CreateState("")
	assert.Empty(state.RedirectURL, "if we provide an empty string, it still shouldn't add a redirect url")

	serialized, err = m.State("foo")
	assert.Nil(err)
	assert.NotEmpty(serialized, "if we provide a valid string it should create a serialized state object")

	deserialized, err := DeserializeState(serialized)
	assert.Nil(err)
	assert.Nil(m.ValidateOAuthState(deserialized))
	assert.Equal("foo", deserialized.RedirectURL)

	m.WithSecret(util.Crypto.MustCreateKey(32))
	serialized, err = m.State("foo")
	assert.Nil(err)
	assert.NotEmpty(serialized, "if we provide a valid string it should create a serialized state object")

	deserialized, err = DeserializeState(serialized)
	assert.Nil(err)
	assert.Nil(m.ValidateOAuthState(deserialized))
	assert.Equal("foo", deserialized.RedirectURL)

	serialized, err = m.State()
	assert.Nil(err)
	assert.NotEmpty(serialized, "state should be issued if we have a secret, even if we don't provide a redirect")

	deserialized, err = DeserializeState(serialized)
	assert.Nil(err)
	assert.Nil(m.ValidateOAuthState(deserialized))
	assert.Empty(deserialized.RedirectURL)
}

func TestManagerValidateNonce(t *testing.T) {
	assert := assert.New(t)

	m := New().WithSecret(util.Crypto.MustCreateKey(32)).WithNonceTimeout(time.Second)

	nonce, err := m.CreateNonce(time.Now().UTC())
	assert.Nil(err)
	assert.Nil(m.ValidateNonce(nonce))

	nonce, err = m.CreateNonce(time.Now().UTC().Add(-time.Minute))
	assert.Nil(err)
	assert.NotNil(m.ValidateNonce(nonce))
}

func TestManagerValidateJWT(t *testing.T) {
	assert := assert.New(t)

	testToken, err := SerializeJWT(util.Crypto.MustCreateKey(32), &JWTPayload{AUD: "client_id"})
	assert.Nil(err)

	jwt, err := DeserializeJWT(testToken)
	assert.Nil(err)
	m := New().WithHostedDomain("blend.com").WithClientID(jwt.Payload.AUD)

	assert.Nil(m.ValidateJWT(jwt))

	m.WithClientID(uuid.V4().String())
	assert.NotNil(m.ValidateJWT(jwt))
}

func TestManagerValidateConfig(t *testing.T) {
	assert := assert.New(t)

	assert.Equal(ErrSecretRequired, New().WithSecret(nil).ValidateConfig())
	assert.Equal(ErrClientIDRequired, New().WithClientID("").ValidateConfig())
	assert.Equal(ErrClientSecretRequired, New().WithClientID("foo").WithClientSecret("").ValidateConfig())
	assert.Equal(ErrRedirectURIRequired, New().WithClientID("foo").WithClientSecret("bar").ValidateConfig())
	assert.Equal(ErrInvalidRedirectURI, New().WithClientID("foo").WithClientSecret("bar").WithRedirectURI(uuid.V4().String()).ValidateConfig())
	assert.Equal(ErrInvalidRedirectURI, New().WithClientID("foo").WithClientSecret("bar").WithRedirectURI("localhost").ValidateConfig())
	assert.Equal(ErrInvalidRedirectURI, New().WithClientID("foo").WithClientSecret("bar").WithRedirectURI("http://").ValidateConfig())
}

func TestManagerFetchProfile(t *testing.T) {
	assert := assert.New(t)
	defer request.ClearMockedResponses()
	request.MockResponseFromString("GET", "https://www.googleapis.com/oauth2/v1/userinfo?access_token=doesnt_matter&alt=json", http.StatusOK, `{"id":"foo", "email":"foo@bar.com"}`)

	m := New().WithClientID("test").WithClientSecret("secret").WithRedirectURI("http://localhost/oauth/finish")

	profile, err := m.FetchProfile("doesnt_matter")
	assert.Nil(err)
	assert.Equal("foo", profile.ID)
	assert.Equal("foo@bar.com", profile.Email)
}

func TestManagerFinish(t *testing.T) {
	assert := assert.New(t)

	m := New().WithClientID("test").WithClientSecret("secret").WithRedirectURI("http://localhost/oauth/finish")

	res, err := m.Finish(&http.Request{URL: &url.URL{}})
	assert.Nil(res)
	assert.Equal(ErrCodeMissing, err)

	res, err = m.Finish(&http.Request{URL: &url.URL{RawQuery: `code=test`}})
	assert.Nil(res)
	assert.Equal(ErrStateMissing, err)
}
