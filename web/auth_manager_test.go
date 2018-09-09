package web

import (
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/util"
	"github.com/blend/go-sdk/uuid"
	"github.com/blend/go-sdk/webutil"
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

	var calledRemoveHandler bool
	removeHandler := am.RemoveHandler()
	am.WithRemoveHandler(func(ctx context.Context, sessionID string, state State) error {
		calledRemoveHandler = true
		if removeHandler == nil {
			return nil
		}
		return removeHandler(ctx, sessionID, state)
	})

	res := NewMockResponseWriter(new(bytes.Buffer))
	r := NewCtx(res, webutil.NewMockRequest("GET", "/"))

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
	assert.False(calledRemoveHandler)

	cookies := ReadSetCookies(res.Header())
	assert.NotEmpty(cookies)
	cookie := cookies[0]
	assert.Equal(am.CookieName(), cookie.Name)
	assert.Equal(am.CookiePath(), cookie.Path)
	assert.Equal(session.SessionID, cookie.Value)
}

func TestAuthManagerLogout(t *testing.T) {
	assert := assert.New(t)

	am := NewLocalAuthManager()

	var calledRemoveHandler bool
	removeHandler := am.RemoveHandler()
	am.WithRemoveHandler(func(ctx context.Context, sessionID string, state State) error {
		calledRemoveHandler = true
		if removeHandler == nil {
			return nil
		}
		return removeHandler(ctx, sessionID, state)
	})

	res := NewMockResponseWriter(new(bytes.Buffer))
	r := NewCtx(res, webutil.NewMockRequest("GET", "/"))

	session, err := am.Login("bailey@blend.com", r)
	assert.Nil(err)
	assert.NotNil(session)

	res = NewMockResponseWriter(new(bytes.Buffer))
	r = NewCtx(res, webutil.NewMockRequestWithCookie("GET", "/", am.CookieName(), session.SessionID))

	assert.Nil(am.Logout(r))
	assert.True(calledRemoveHandler)

	cookies := ReadSetCookies(res.Header())
	assert.NotEmpty(cookies)
	cookie := cookies[0]
	assert.Equal(am.CookieName(), cookie.Name)
	assert.Equal(am.CookiePath(), cookie.Path)
	assert.NotEqual(session.SessionID, cookie.Value, "we should randomize the session cookie on logout")
	assert.True(time.Now().UTC().After(cookie.Expires))
}

func TestAuthManagerVerifySessionParsed(t *testing.T) {
	assert := assert.New(t)

	am := NewLocalAuthManager()

	var calledParseHandler bool
	am.WithParseSessionValueHandler(func(ctx context.Context, sessionID string, state State) (*Session, error) {
		calledParseHandler = true
		return &Session{UserID: uuid.V4().String(), SessionID: sessionID}, nil
	})

	var calledFetchHandler bool
	am.WithFetchHandler(func(ctx context.Context, sessionID string, state State) (*Session, error) {
		calledFetchHandler = true
		return nil, nil
	})

	var calledValidateHandler bool
	am.WithValidateHandler(func(ctx context.Context, session *Session, state State) error {
		calledValidateHandler = true
		return nil
	})

	assert.NotNil(am.ParseSessionValueHandler())
	assert.NotNil(am.FetchHandler())
	assert.NotNil(am.ValidateHandler())

	r := NewCtx(NewMockResponseWriter(new(bytes.Buffer)), webutil.NewMockRequestWithCookie("GET", "/", am.CookieName(), NewSessionID()))
	session, err := am.VerifySession(r)
	assert.Nil(err)
	assert.NotNil(session)
	assert.True(session.ExpiresUTC.IsZero())
	assert.True(calledParseHandler)
	assert.False(calledFetchHandler)
	assert.True(calledValidateHandler)
}

func TestAuthManagerVerifySessionFetched(t *testing.T) {
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

	r := NewCtx(NewMockResponseWriter(new(bytes.Buffer)), webutil.NewMockRequest("GET", "/"))
	session, err := am.Login("bailey@blend.com", r)
	assert.Nil(err)
	assert.NotNil(session)
	assert.False(calledFetchHandler)
	assert.False(calledValidateHandler)

	r = NewCtx(NewMockResponseWriter(new(bytes.Buffer)), webutil.NewMockRequestWithCookie("GET", "/", am.CookieName(), session.SessionID))
	session, err = am.VerifySession(r)
	assert.Nil(err)
	assert.NotNil(session)
	assert.True(calledFetchHandler)
	assert.True(calledValidateHandler)
}

func TestAuthManagerVerifySessionUnset(t *testing.T) {
	assert := assert.New(t)

	am := NewLocalAuthManager()

	r := NewCtx(NewMockResponseWriter(new(bytes.Buffer)), webutil.NewMockRequest("GET", "/"))

	session, err := am.VerifySession(r)
	assert.Nil(err)
	assert.Nil(session)
}

func TestAuthManagerVerifySessionExpired(t *testing.T) {
	assert := assert.New(t)

	am := NewLocalAuthManager()
	am.WithSessionTimeoutProvider(nil)
	am.WithParseSessionValueHandler(func(ctx context.Context, sessionID string, state State) (*Session, error) {
		return &Session{UserID: uuid.V4().String(), SessionID: sessionID, ExpiresUTC: time.Now().UTC().Add(-time.Hour)}, nil
	})

	var calledValidateHandler bool
	am.WithValidateHandler(func(ctx context.Context, session *Session, state State) error {
		calledValidateHandler = true
		return nil
	})

	res := NewMockResponseWriter(new(bytes.Buffer))
	r := NewCtx(res, webutil.NewMockRequestWithCookie("GET", "/", am.CookieName(), NewSessionID()))
	session, err := am.VerifySession(r)
	assert.Nil(err)
	assert.Nil(session)
	assert.False(calledValidateHandler)

	// assert the cookie is expired ...
	cookies := ReadSetCookies(res.Header())
	assert.NotEmpty(cookies)

	cookie := cookies[0]
	assert.Equal(am.CookieName(), cookie.Name)
	assert.Equal(am.CookiePath(), cookie.Path)
	assert.True(cookie.Expires.Before(time.Now().UTC()), "the cookie should be expired")
}
