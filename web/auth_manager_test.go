package web

import (
	"bytes"
	"testing"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/util"
)

func TestNewJWTAuthManager(t *testing.T) {
	assert := assert.New(t)

	am := NewJWTAuthManager(util.Crypto.MustCreateKey(64))
	assert.NotNil(am.SessionTimeoutProvider(), "must set a session timeout provider for a jwt manager")
}

func TestAuthManagerLogin(t *testing.T) {
	assert := assert.New(t)

	am := NewLocalAuthManager()

	r := NewCtx(NewMockResponseWriter(new(bytes.Buffer)), NewMockRequest("GET", "/"), nil, nil)

	session, err := am.Login("bailey@blend.com", r)
	assert.Nil(err)
	assert.NotNil(session)
	assert.NotEmpty(session.SessionID)
	assert.NotEmpty(session.RemoteAddr)
	assert.NotEmpty(session.UserAgent)
	assert.Equal("bailey@blend.com", session.UserID)
	assert.True(session.ExpiresUTC.IsZero())
}

func TestAuthManagerVerifySession(t *testing.T) {
	assert := assert.New(t)

	am := NewLocalAuthManager()

	r := NewCtx(NewMockResponseWriter(new(bytes.Buffer)), NewMockRequest("GET", "/"), nil, nil)

	session, err := am.VerifySession(r)
	assert.Nil(err)
	assert.Nil(session)

	r = NewCtx(NewMockResponseWriter(new(bytes.Buffer)), NewMockRequest("GET", "/"), nil, nil)
	session, err = am.Login("bailey@blend.com", r)
	assert.Nil(err)
	assert.NotNil(session)

	r = NewCtx(NewMockResponseWriter(new(bytes.Buffer)), NewMockRequestWithCookie("GET", "/", am.CookieName(), session.SessionID), nil, nil)
	session, err = am.VerifySession(r)
	assert.Nil(err)
	assert.NotNil(session)
}
