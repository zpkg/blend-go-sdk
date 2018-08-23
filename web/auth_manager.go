package web

import (
	"net/url"
	"time"

	"github.com/blend/go-sdk/logger"
)

const (
	// AuthManagerModeJWT is the jwt auth mode.
	AuthManagerModeJWT = "jwt"
	// AuthManagerModeServer is the server managed auth mode.
	AuthManagerModeServer = "server"
	// AuthManagerModeLocal is the local cache auth mode.
	AuthManagerModeLocal = "cached"
)

// NewAuthManagerFromConfig returns a new auth manager from a given config.
func NewAuthManagerFromConfig(cfg *Config) (manager *AuthManager) {
	switch cfg.GetAuthManagerMode() {
	case AuthManagerModeJWT:
		manager = NewJWTAuthManager(cfg.GetAuthSecret())
	case AuthManagerModeServer:
		manager = NewServerAuthManager()
	case AuthManagerModeLocal:
		fallthrough
	default:
		manager = NewLocalAuthManager()
	}

	return manager.WithCookieHTTPSOnly(cfg.GetCookieHTTPSOnly()).
		WithCookieName(cfg.GetCookieName()).
		WithCookiePath(cfg.GetCookiePath()).
		WithSessionTimeoutProvider(SessionTimeoutProvider(cfg.GetSessionTimeoutIsAbsolute(), cfg.GetSessionTimeout()))
}

// NewLocalAuthManager returns a new locally cached session manager.
// It saves sessions to a local store.
func NewLocalAuthManager() *AuthManager {
	cache := NewLocalSessionCache()
	return &AuthManager{
		persistHandler: cache.PersistHandler,
		fetchHandler:   cache.FetchHandler,
		removeHandler:  cache.RemoveHandler,
		cookieName:     DefaultCookieName,
		cookiePath:     DefaultCookiePath,
	}
}

// NewJWTAuthManager returns a new jwt session manager.
// It issues JWT tokens to identify users.
func NewJWTAuthManager(key []byte) *AuthManager {
	jwtm := NewJWTManagerForKey(key)
	return &AuthManager{
		serializeSessionValueHandler: jwtm.SerializeSessionValueHandler,
		parseSessionValueHandler:     jwtm.ParseSessionValueHandler,
		cookieName:                   DefaultCookieName,
		cookiePath:                   DefaultCookiePath,
	}
}

// NewServerAuthManager returns a new server auth manager.
// You should set the `FetchHandler`, the `PersistHandler` and the `RemoveHandler`.
func NewServerAuthManager() *AuthManager {
	return &AuthManager{
		cookieName: DefaultCookieName,
		cookiePath: DefaultCookiePath,
	}
}

// AuthManager is a manager for sessions.
type AuthManager struct {
	// these generally apply to jwt mode.
	serializeSessionValueHandler func(*Session, State) (string, error)
	parseSessionValueHandler     func(string, State) (*Session, error) // should provide the server tracked session id or should return the jwt.

	// these generally apply to server or local modes.
	persistHandler func(*Session, State) error
	fetchHandler   func(sessionID string, state State) (*Session, error)
	removeHandler  func(sessionID string, state State) error

	// these generally apply to any mode.
	validateHandler          func(*Session, State) error
	sessionTimeoutProvider   func(*Session) time.Time
	loginRedirectHandler     func(*Ctx) *url.URL
	postLoginRedirectHandler func(*Ctx) *url.URL

	cookieName      string
	cookiePath      string
	cookieHTTPSOnly bool
}

// --------------------------------------------------------------------------------
// Methods
// --------------------------------------------------------------------------------

// Login logs a userID in.
func (am *AuthManager) Login(userID string, ctx *Ctx) (session *Session, err error) {
	// create a new session value
	sessionValue := am.createSessionID()
	// userID and sessionID are required
	session = NewSession(userID, sessionValue)
	session.ExpiresUTC = am.GenerateSessionTimeout(ctx)
	session.UserAgent = logger.GetUserAgent(ctx.request)
	session.RemoteAddr = logger.GetRemoteAddr(ctx.request)

	// call the perist handler if one's been provided
	if am.persistHandler != nil {
		err = am.persistHandler(session, ctx.state)
		if err != nil {
			return nil, err
		}
	}

	// if we're in jwt mode, serialize the jwt.
	if am.serializeSessionValueHandler != nil {
		sessionValue, err = am.serializeSessionValueHandler(session, ctx.state)
		if err != nil {
			return nil, err
		}
	}

	// inject cookies into the response
	am.injectCookie(ctx, am.CookieName(), sessionValue, session.ExpiresUTC)
	return session, nil
}

// Logout unauthenticates a session.
func (am *AuthManager) Logout(ctx *Ctx) error {
	sessionValue := am.readSessionValue(ctx)

	// issue the expiration cookies to the response
	ctx.ExpireCookie(am.CookieName(), am.CookiePath())
	// nil out the current session in the ctx
	ctx.WithSession(nil)

	// call the remove handler if one has been provided
	if am.removeHandler != nil {
		return am.removeHandler(sessionValue, ctx.state)
	}
	return nil
}

// VerifySession checks a sessionID to see if it's valid.
// It also handles updating a rolling expiry.
func (am *AuthManager) VerifySession(ctx *Ctx) (*Session, error) {
	var err error
	// pull the sessionID off the request
	sessionValue := am.readSessionValue(ctx)

	// validate the sessionValue isn't unset or crazy long.
	if err = am.sanityCheckSessionValue(sessionValue); err != nil {
		return nil, err
	}
	var session *Session
	// if we have a separate step to parse the sesion value
	// (i.e. jwt mode) do that now.
	if am.parseSessionValueHandler != nil {
		session, err = am.parseSessionValueHandler(sessionValue, ctx.state)
		if err != nil {
			return nil, err
		}
	} else if am.fetchHandler != nil { // if we're in server tracked mode, pull it from whatever backing store we use.
		session, err = am.fetchHandler(sessionValue, ctx.state)
		if err != nil {
			return nil, err
		}
	}

	// if the session is invalid, expire the cookie(s)
	if session == nil || session.IsZero() || session.IsExpired() {
		ctx.ExpireCookie(am.CookieName(), am.CookiePath())
		// if we have a remove handler and the sessionID is set
		if am.removeHandler != nil {
			err = am.removeHandler(sessionValue, ctx.state)
			if err != nil {
				return nil, err
			}
		}

		// exit out, the session is bad
		return nil, nil
	}

	// call a custom validate handler if one's been provided.
	if am.validateHandler != nil {
		err = am.validateHandler(session, ctx.state)
		if err != nil {
			return nil, err
		}
	}

	// check if we need to do a rolling expiry update
	// note this will be explicitly false by default
	// as we use absolte expiry by default.
	if am.shouldUpdateSessionExpiry() {
		session.ExpiresUTC = am.GenerateSessionTimeout(ctx)
		if am.persistHandler != nil {
			err = am.persistHandler(session, ctx.state)
			if err != nil {
				return nil, err
			}
		}

		am.injectCookie(ctx, am.CookieName(), sessionValue, session.ExpiresUTC)
	}

	return session, nil
}

// LoginRedirect returns a redirect result for when auth fails and you need to
// send the user to a login page.
func (am *AuthManager) LoginRedirect(ctx *Ctx) Result {
	if am.loginRedirectHandler != nil {
		redirectTo := am.loginRedirectHandler(ctx)
		if redirectTo != nil {
			return ctx.Redirect(redirectTo.String())
		}
	}
	return ctx.DefaultResultProvider().NotAuthorized()
}

// PostLoginRedirect returns a redirect result for when auth fails and you need to
// send the user to a login page.
func (am *AuthManager) PostLoginRedirect(ctx *Ctx) Result {
	if am.postLoginRedirectHandler != nil {
		redirectTo := am.postLoginRedirectHandler(ctx)
		if redirectTo != nil {
			return ctx.Redirect(redirectTo.String())
		}
	}
	// the default authed redirect is the root.
	return ctx.RedirectWithMethod("GET", "/")
}

// --------------------------------------------------------------------------------
// Properties
// --------------------------------------------------------------------------------

// WithSessionTimeoutProvider sets the session timeout provider.
func (am *AuthManager) WithSessionTimeoutProvider(timeoutProvider func(*Session) time.Time) *AuthManager {
	am.sessionTimeoutProvider = timeoutProvider
	return am
}

// SessionTimeoutProvider returns the session timeout provider.
func (am *AuthManager) SessionTimeoutProvider() func(*Session) time.Time {
	return am.sessionTimeoutProvider
}

// WithCookieHTTPSOnly sets if we should issue cookies with the HTTPS flag on.
func (am *AuthManager) WithCookieHTTPSOnly(isHTTPSOnly bool) *AuthManager {
	am.cookieHTTPSOnly = isHTTPSOnly
	return am
}

// CookiesHTTPSOnly returns if the cookie is for only https connections.
func (am *AuthManager) CookiesHTTPSOnly() bool {
	return am.cookieHTTPSOnly
}

// WithCookieName sets the cookie name.
func (am *AuthManager) WithCookieName(paramName string) *AuthManager {
	am.cookieName = paramName
	return am
}

// CookieName returns the session param name.
func (am *AuthManager) CookieName() string {
	return am.cookieName
}

// WithCookiePath sets the cookie path.
func (am *AuthManager) WithCookiePath(path string) *AuthManager {
	am.cookiePath = path
	return am
}

// CookiePath returns the session param path.
func (am *AuthManager) CookiePath() string {
	if len(am.cookiePath) == 0 {
		return DefaultCookiePath
	}
	return am.cookiePath
}

// WithPersistHandler sets the persist handler.
func (am *AuthManager) WithPersistHandler(handler func(*Session, State) error) *AuthManager {
	am.persistHandler = handler
	return am
}

// PersistHandler returns the persist handler.
func (am *AuthManager) PersistHandler() func(*Session, State) error {
	return am.persistHandler
}

// WithFetchHandler sets the fetch handler.
func (am *AuthManager) WithFetchHandler(handler func(sessionID string, state State) (*Session, error)) *AuthManager {
	am.fetchHandler = handler
	return am
}

// FetchHandler returns the fetch handler.
// It is used in `VerifySession` to satisfy session cache misses.
func (am *AuthManager) FetchHandler() func(sessionID string, state State) (*Session, error) {
	return am.fetchHandler
}

// WithRemoveHandler sets the remove handler.
func (am *AuthManager) WithRemoveHandler(handler func(sessionID string, state State) error) *AuthManager {
	am.removeHandler = handler
	return am
}

// RemoveHandler returns the remove handler.
// It is used in validate session if the session is found to be invalid.
func (am *AuthManager) RemoveHandler() func(sessionID string, state State) error {
	return am.removeHandler
}

// WithValidateHandler sets the validate handler.
func (am *AuthManager) WithValidateHandler(handler func(*Session, State) error) *AuthManager {
	am.validateHandler = handler
	return am
}

// ValidateHandler returns the validate handler.
func (am *AuthManager) ValidateHandler() func(*Session, State) error {
	return am.validateHandler
}

// WithLoginRedirectHandler sets the login redirect handler.
func (am *AuthManager) WithLoginRedirectHandler(handler func(*Ctx) *url.URL) *AuthManager {
	am.loginRedirectHandler = handler
	return am
}

// LoginRedirectHandler returns the login redirect handler.
func (am *AuthManager) LoginRedirectHandler() func(*Ctx) *url.URL {
	return am.loginRedirectHandler
}

// WithPostLoginRedirectHandler sets the post login redirect handler.
func (am *AuthManager) WithPostLoginRedirectHandler(handler func(*Ctx) *url.URL) *AuthManager {
	am.postLoginRedirectHandler = handler
	return am
}

// PostLoginRedirectHandler returns the redirect handler for login complete.
func (am *AuthManager) PostLoginRedirectHandler() func(*Ctx) *url.URL {
	return am.postLoginRedirectHandler
}

// --------------------------------------------------------------------------------
// Utility Methods
// --------------------------------------------------------------------------------

// GenerateSessionTimeout returns the absolute time the cookie would expire.
func (am *AuthManager) GenerateSessionTimeout(context *Ctx) (output time.Time) {
	if am.sessionTimeoutProvider != nil {
		output = am.sessionTimeoutProvider(context.Session())
	}
	return
}

func (am AuthManager) shouldUpdateSessionExpiry() bool {
	return am.sessionTimeoutProvider != nil
}

// CreateSessionID creates a new session id.
func (am AuthManager) createSessionID() string {
	return NewSessionID()
}

// InjectCookie injects a session cookie into the context.
func (am *AuthManager) injectCookie(ctx *Ctx, name, value string, expire time.Time) {
	ctx.WriteNewCookie(name, value, expire, am.CookiePath(), am.CookiesHTTPSOnly())
}

// readParam reads a param from a given request context from either the cookies or headers.
func (am *AuthManager) readParam(name string, ctx *Ctx) (output string) {
	if cookie := ctx.GetCookie(name); cookie != nil {
		output = cookie.Value
	}
	return
}

// ReadSessionID reads a session id from a given request context.
func (am *AuthManager) readSessionValue(ctx *Ctx) string {
	return am.readParam(am.CookieName(), ctx)
}

// ValidateSessionID verifies a session id.
func (am *AuthManager) sanityCheckSessionValue(sessionID string) error {
	if len(sessionID) == 0 {
		return ErrSessionIDEmpty
	}
	return nil
}
