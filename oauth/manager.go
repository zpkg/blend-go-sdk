package oauth

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/blend/go-sdk/exception"
	"github.com/blend/go-sdk/logger"
	"github.com/blend/go-sdk/request"
	"github.com/blend/go-sdk/util"
	"github.com/blend/go-sdk/uuid"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// New returns a new manager.
// By default it will error if you try and validate a profile.
// You must either enable `SkipDomainvalidation` or provide valid domains.
func New() *Manager {
	return &Manager{}
}

// Must is a helper for handling NewFromEnv() and NewFromConfig().
func Must(m *Manager, err error) *Manager {
	if err != nil {
		panic(err)
	}
	return m
}

// NewFromEnv returns a new manager from the environment.
func NewFromEnv() (*Manager, error) {
	return NewFromConfig(NewConfigFromEnv())
}

// NewFromConfig returns a new oauth manager from a config.
func NewFromConfig(cfg *Config) (*Manager, error) {
	secret, err := cfg.GetSecret()
	if err != nil {
		return nil, err
	}
	return &Manager{
		secret:       secret,
		redirectURI:  cfg.GetRedirectURI(),
		hostedDomain: cfg.GetHostedDomain(),
		scopes:       cfg.GetScopes(),
		clientID:     cfg.GetClientID(),
		clientSecret: cfg.GetClientSecret(),
	}, nil
}

// Manager is the oauth manager.
type Manager struct {
	secret       []byte
	scopes       []string
	redirectURI  string
	hostedDomain string
	clientID     string
	clientSecret string
}

func (m *Manager) conf(r *http.Request) *oauth2.Config {
	return &oauth2.Config{
		ClientID:     m.clientID,
		ClientSecret: m.clientSecret,
		RedirectURL:  m.getRedirectURI(r),
		Scopes:       m.scopes,
		Endpoint:     google.Endpoint,
	}
}

func (m *Manager) getRedirectURI(r *http.Request) string {
	if util.String.HasPrefixCaseInsensitive(m.redirectURI, "https://") ||
		util.String.HasPrefixCaseInsensitive(m.redirectURI, "http://") ||
		util.String.HasPrefixCaseInsensitive(m.redirectURI, "spdy://") {
		return m.redirectURI
	}

	requestURI := &url.URL{
		Scheme: logger.GetProto(r),
		Host:   logger.GetHost(r),
		Path:   m.redirectURI,
	}
	return requestURI.String()
}

// OAuthURL is the auth url for google with a given clientID.
// This is typically the link that a user will click on to start the auth process.
func (m *Manager) OAuthURL(r *http.Request, redirect ...string) (oauthURL string, err error) {
	var state string
	state, err = SerializeState(m.CreateState(redirect...))
	if err != nil {
		return
	}

	var opts []oauth2.AuthCodeOption
	if len(m.hostedDomain) > 0 {
		opts = append(opts, oauth2.SetAuthURLParam("hd", m.hostedDomain))
	}
	oauthURL = m.conf(r).AuthCodeURL(state, opts...)
	return
}

// Finish processes the returned code, exchanging for an access token, and fetches the user profile.
func (m *Manager) Finish(r *http.Request) (*Result, error) {
	var result Result
	var err error

	// grab the code off the request.
	code := r.URL.Query().Get("code")
	if len(code) == 0 {
		return nil, ErrCodeMissing
	}

	// fetch the state
	state := r.URL.Query().Get("state")
	if len(state) > 0 {
		deserialized, err := DeserializeState(state)
		if err != nil {
			return nil, err
		}
		result.State = deserialized
	}

	err = m.ValidateState(result.State)
	if err != nil {
		return nil, err
	}

	// Handle the exchange code to initiate a transport.
	tok, err := m.conf(r).Exchange(oauth2.NoContext, code)
	if err != nil {
		return nil, exception.New(ErrFailedCodeExchange).WithMessagef("inner: %+v", err)
	}
	result.Response.AccessToken = tok.AccessToken
	result.Response.TokenType = tok.TokenType
	result.Response.RefreshToken = tok.RefreshToken
	result.Response.Expiry = tok.Expiry

	prof, err := m.FetchProfile(tok.AccessToken)
	if err != nil {
		return nil, err
	}
	result.Profile = prof
	return &result, nil
}

// FetchProfile gets a google profile for an access token.
func (m *Manager) FetchProfile(accessToken string) (profile Profile, err error) {
	req, err := request.New().AsGet().
		WithRawURL("https://www.googleapis.com/oauth2/v1/userinfo")

	contents, meta, err := req.
		WithQueryString("alt", "json").
		WithQueryString("access_token", accessToken).
		WithMockProvider(request.MockedResponseInjector).
		BytesWithMeta()

	if err != nil {
		return
	}
	if meta.StatusCode > 299 {
		err = exception.New(ErrGoogleResponseStatus).WithMessagef("status code: %d, response: %s", meta.StatusCode, string(contents))
		return
	}
	if err = json.Unmarshal(contents, &profile); err != nil {
		err = exception.New(ErrProfileJSONUnmarshal).WithMessagef("inner: %v", err)
		return
	}
	return
}

// CreateState creates auth state.
func (m *Manager) CreateState(redirect ...string) (state State) {
	if len(m.secret) > 0 {
		state.Token = uuid.V4().String()
		state.SecureToken = m.hash(state.Token)
	}

	if len(redirect) > 0 && len(redirect[0]) > 0 {
		state.RedirectURL = redirect[0]
	}
	return
}

// --------------------------------------------------------------------------------
// Validation Helpers
// --------------------------------------------------------------------------------

// ValidateState validates oauth state.
func (m *Manager) ValidateState(state State) error {
	if len(m.secret) > 0 {
		expected := m.hash(state.Token)
		actual := state.SecureToken
		if !hmac.Equal([]byte(expected), []byte(actual)) {
			return ErrInvalidAntiforgeryToken
		}
	}
	return nil
}

// ValidateProfile validates a profile.
func (m *Manager) ValidateProfile(p *Profile) error {
	if len(m.HostedDomain()) == 0 {
		return nil
	}

	workingDomain := m.hostedDomain
	if !strings.HasPrefix(workingDomain, "@") {
		workingDomain = fmt.Sprintf("@%s", workingDomain)
	}
	if !util.String.HasSuffixCaseInsensitive(p.Email, workingDomain) {
		return ErrInvalidHostedDomain
	}
	return nil
}

// --------------------------------------------------------------------------------
// Properties
// --------------------------------------------------------------------------------

// WithSecret sets the secret used to create state tokens.
func (m *Manager) WithSecret(secret []byte) *Manager {
	m.secret = secret
	return m
}

// Secret returns a property
func (m *Manager) Secret() []byte {
	return m.secret
}

// WithRedirectURI sets the return url.
func (m *Manager) WithRedirectURI(redirectURI string) *Manager {
	m.redirectURI = redirectURI
	return m
}

// RedirectURI returns a property.
func (m *Manager) RedirectURI() string {
	return m.redirectURI
}

// WithHostedDomain returns the hosted domain.
func (m *Manager) WithHostedDomain(hostedDomain string) *Manager {
	m.hostedDomain = hostedDomain
	return m
}

// HostedDomain returns the hosted domain.
func (m *Manager) HostedDomain() string {
	return m.hostedDomain
}

// WithClientID sets the client id.
func (m *Manager) WithClientID(clientID string) *Manager {
	m.clientID = clientID
	return m
}

// WithScopes sets the oauth scopes.
func (m *Manager) WithScopes(scopes ...string) *Manager {
	m.scopes = scopes
	return m
}

// Scopes returns the oauth scopes.
func (m *Manager) Scopes() []string {
	return m.scopes
}

// ClientID returns a property.
func (m *Manager) ClientID() string {
	return m.clientID
}

// WithClientSecret sets the client id.
func (m *Manager) WithClientSecret(clientSecret string) *Manager {
	m.clientSecret = clientSecret
	return m
}

// ClientSecret returns a client secret.
func (m *Manager) ClientSecret() string {
	return m.clientSecret
}

// --------------------------------------------------------------------------------
// internal helpers
// --------------------------------------------------------------------------------

func (m *Manager) hash(plaintext string) string {
	return base64.URLEncoding.EncodeToString(m.hmac([]byte(plaintext)))
}

// hmac hashes data with the given key.
func (m *Manager) hmac(plainText []byte) []byte {
	mac := hmac.New(sha512.New, m.secret)
	mac.Write([]byte(plainText))
	return mac.Sum(nil)
}
