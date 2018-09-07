package web

import (
	"bytes"
	"context"
	"fmt"
	"strings"
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

	var calledPersistHandler bool
	persistHandler := am.PersistHandler()
	am.WithPersistHandler(func(ctx context.Context, session *Session, state State) error {
		calledPersistHandler = true
		if persistHandler == nil {
			return nil
		}
		return persistHandler(ctx, session, state)
	})

	var calledSerializeHandler bool
	serializeHandler := am.SerializeSessionValueHandler()
	am.WithSerializeSessionValueHandler(func(ctx context.Context, session *Session, state State) (string, error) {
		calledSerializeHandler = true
		if serializeHandler == nil {
			return session.SessionID, nil
		}
		return serializeHandler(ctx, session, state)
	})

	res := NewMockResponseWriter(new(bytes.Buffer))
	r := NewCtx(res, NewMockRequest("GET", "/"), nil, nil)

	session, err := am.Login("bailey@blend.com", r)
	assert.Nil(err)
	assert.NotNil(session)
	assert.NotEmpty(session.SessionID)
	assert.NotEmpty(session.RemoteAddr)
	assert.NotEmpty(session.UserAgent)
	assert.Equal("bailey@blend.com", session.UserID)
	assert.True(session.ExpiresUTC.IsZero())
	assert.True(calledPersistHandler)
	assert.True(calledSerializeHandler)
	cookie := res.Header().Get("Set-Cookie")
	assert.True(strings.HasPrefix(cookie, fmt.Sprintf("%s=%s", am.CookieName(), session.SessionID)))
}

func TestAuthManagerVerifySession(t *testing.T) {
	assert := assert.New(t)

	am := NewLocalAuthManager()

	var calledFetchHandler bool
	fetchHandler := am.FetchHandler()
	am.WithFetchHandler(func(ctx context.Context, sessionID string, state State) (*Session, error) {
		calledFetchHandler = true
		if fetchHandler == nil {
			return nil, nil
		}
		return fetchHandler(ctx, sessionID, state)
	})

	var calledValidateHandler bool
	validateHandler := am.ValidateHandler()
	am.WithValidateHandler(func(ctx context.Context, session *Session, state State) error {
		calledValidateHandler = true
		if validateHandler == nil {
			return nil
		}
		return validateHandler(ctx, session, state)
	})

	r := NewCtx(NewMockResponseWriter(new(bytes.Buffer)), NewMockRequest("GET", "/"), nil, nil)
	session, err := am.Login("bailey@blend.com", r)
	assert.Nil(err)
	assert.NotNil(session)
	assert.False(calledFetchHandler)

	r = NewCtx(NewMockResponseWriter(new(bytes.Buffer)), NewMockRequestWithCookie("GET", "/", am.CookieName(), session.SessionID), nil, nil)
	session, err = am.VerifySession(r)
	assert.Nil(err)
	assert.NotNil(session)
	assert.True(calledFetchHandler)
	assert.True(calledValidateHandler)
}

func TestAuthManagerVerifySessionUnset(t *testing.T) {
	assert := assert.New(t)

	am := NewLocalAuthManager()

	r := NewCtx(NewMockResponseWriter(new(bytes.Buffer)), NewMockRequest("GET", "/"), nil, nil)

	session, err := am.VerifySession(r)
	assert.Nil(err)
	assert.Nil(session)
}

func TestAuthManagerVerifySessionExpired(t *testing.T) {
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
