package web

import (
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/blend/go-sdk/configutil"
	"github.com/blend/go-sdk/env"
	"github.com/blend/go-sdk/webutil"
)

// Config is an object used to set up a web app.
type Config struct {
	Port                      int32             `json:"port,omitempty" yaml:"port,omitempty" env:"PORT"`
	BindAddr                  string            `json:"bindAddr,omitempty" yaml:"bindAddr,omitempty" env:"BIND_ADDR"`
	BaseURL                   string            `json:"baseURL,omitempty" yaml:"baseURL,omitempty" env:"BASE_URL"`
	SkipRedirectTrailingSlash *bool             `json:"skipRedirectTrailingSlash,omitempty" yaml:"skipRedirectTrailingSlash,omitempty"`
	HandleOptions             *bool             `json:"handleOptions,omitempty" yaml:"handleOptions,omitempty"`
	HandleMethodNotAllowed    *bool             `json:"handleMethodNotAllowed,omitempty" yaml:"handleMethodNotAllowed,omitempty"`
	RecoverPanics             *bool             `json:"recoverPanics,omitempty" yaml:"recoverPanics,omitempty"`
	AuthManagerMode           string            `json:"authManagerMode" yaml:"authManagerMode"`
	AuthSecret                string            `json:"authSecret" yaml:"authSecret" env:"AUTH_SECRET"`
	SessionTimeout            time.Duration     `json:"sessionTimeout,omitempty" yaml:"sessionTimeout,omitempty" env:"SESSION_TIMEOUT"`
	SessionTimeoutIsAbsolute  *bool             `json:"sessionTimeoutIsAbsolute,omitempty" yaml:"sessionTimeoutIsAbsolute,omitempty" env:"SESSION_TIMEOUT_ABSOLUTE"`
	CookieHTTPSOnly           *bool             `json:"cookieHTTPSOnly,omitempty" yaml:"cookieHTTPSOnly,omitempty" env:"COOKIE_HTTPS_ONLY"`
	CookieName                string            `json:"cookieName,omitempty" yaml:"cookieName,omitempty" env:"COOKIE_NAME"`
	CookiePath                string            `json:"cookiePath,omitempty" yaml:"cookiePath,omitempty" env:"COOKIE_PATH"`
	DefaultHeaders            map[string]string `json:"defaultHeaders,omitempty" yaml:"defaultHeaders,omitempty"`
	MaxHeaderBytes            int               `json:"maxHeaderBytes,omitempty" yaml:"maxHeaderBytes,omitempty" env:"MAX_HEADER_BYTES"`
	ReadTimeout               time.Duration     `json:"readTimeout,omitempty" yaml:"readTimeout,omitempty" env:"READ_HEADER_TIMEOUT"`
	ReadHeaderTimeout         time.Duration     `json:"readHeaderTimeout,omitempty" yaml:"readHeaderTimeout,omitempty" env:"READ_HEADER_TIMEOUT"`
	WriteTimeout              time.Duration     `json:"writeTimeout,omitempty" yaml:"writeTimeout,omitempty" env:"WRITE_TIMEOUT"`
	IdleTimeout               time.Duration     `json:"idleTimeout,omitempty" yaml:"idleTimeout,omitempty" env:"IDLE_TIMEOUT"`
	ShutdownGracePeriod       time.Duration     `json:"shutdownGracePeriod" yaml:"shutdownGracePeriod" env:"SHUTDOWN_GRACE_PERIOD"`

	Views ViewCacheConfig `json:"views,omitempty" yaml:"views,omitempty"`
}

// Resolve resolves the config from other sources.
func (c *Config) Resolve() error {
	return env.Env().ReadInto(c)
}

// BindAddrOrDefault returns the bind address or a default.
func (c Config) BindAddrOrDefault(defaults ...string) string {
	if len(c.BindAddr) > 0 {
		return c.BindAddr
	}
	if c.Port > 0 {
		return fmt.Sprintf(":%d", c.Port)
	}
	if len(defaults) > 0 {
		return defaults[0]
	}
	return DefaultBindAddr
}

// PortOrDefault returns the int32 port for a given config.
// This is useful in things like kubernetes pod templates.
// If the config .Port is unset, it will parse the .BindAddr,
// or the DefaultBindAddr for the port number.
func (c Config) PortOrDefault() int32 {
	if c.Port > 0 {
		return c.Port
	}
	if len(c.BindAddr) > 0 {
		return webutil.PortFromBindAddr(c.BindAddr)
	}
	return webutil.PortFromBindAddr(DefaultBindAddr)
}

// BaseURLOrDefault gets the base url for the app or a default.
func (c Config) BaseURLOrDefault() string {
	return c.BaseURL
}

// SkipRedirectTrailingSlashOrDefault returns if we should skip automatically redirecting for a missing trailing slash.
func (c Config) SkipRedirectTrailingSlashOrDefault() bool {
	return configutil.CoalesceBool(c.SkipRedirectTrailingSlash, DefaultSkipRedirectTrailingSlash)
}

// HandleOptionsOrDefault returns if we should handle OPTIONS verb requests.
func (c Config) HandleOptionsOrDefault() bool {
	return configutil.CoalesceBool(c.HandleOptions, DefaultHandleOptions)
}

// HandleMethodNotAllowedOrDefault returns if we should handle method not allowed results.
func (c Config) HandleMethodNotAllowedOrDefault() bool {
	return configutil.CoalesceBool(c.HandleMethodNotAllowed, DefaultHandleMethodNotAllowed)
}

// RecoverPanicsOrDefault returns if we should recover panics or not.
func (c Config) RecoverPanicsOrDefault() bool {
	return configutil.CoalesceBool(c.RecoverPanics, DefaultRecoverPanics)
}

// BaseURLIsSecureScheme returns if the base url starts with a secure scheme.
func (c Config) BaseURLIsSecureScheme() bool {
	if c.BaseURL == "" {
		return false
	}
	return strings.HasPrefix(strings.ToLower(c.BaseURL), SchemeHTTPS) || strings.HasPrefix(strings.ToLower(c.BaseURL), SchemeSPDY)
}

// AuthManagerModeOrDefault returns the auth manager mode.
func (c Config) AuthManagerModeOrDefault() AuthManagerMode {
	return AuthManagerMode(configutil.CoalesceString(c.AuthManagerMode, string(AuthManagerModeRemote)))
}

// MustAuthSecret returns the auth secret and panics if there is an error decoding it.
func (c Config) MustAuthSecret() []byte {
	decoded, err := base64.StdEncoding.DecodeString(c.AuthSecret)
	if err != nil {
		panic(err)
	}
	return decoded
}

// SessionTimeoutOrDefault returns a property or a default.
func (c Config) SessionTimeoutOrDefault() time.Duration {
	return configutil.CoalesceDuration(c.SessionTimeout, DefaultSessionTimeout)
}

// SessionTimeoutIsAbsoluteOrDefault returns a property or a default.
func (c Config) SessionTimeoutIsAbsoluteOrDefault() bool {
	return configutil.CoalesceBool(c.SessionTimeoutIsAbsolute, DefaultSessionTimeoutIsAbsolute)
}

// CookieHTTPSOnlyOrDefault returns a property or a default.
func (c Config) CookieHTTPSOnlyOrDefault() bool {
	return configutil.CoalesceBool(c.CookieHTTPSOnly, true)
}

// CookieNameOrDefault returns a property or a default.
func (c Config) CookieNameOrDefault() string {
	return configutil.CoalesceString(c.CookieName, DefaultCookieName)
}

// CookiePathOrDefault returns a property or a default.
func (c Config) CookiePathOrDefault() string {
	return configutil.CoalesceString(c.CookiePath, DefaultCookiePath)
}

// MaxHeaderBytesOrDefault returns the maximum header size in bytes or a default.
func (c Config) MaxHeaderBytesOrDefault() int {
	return configutil.CoalesceInt(c.MaxHeaderBytes, DefaultMaxHeaderBytes)
}

// ReadTimeoutOrDefault gets a property.
func (c Config) ReadTimeoutOrDefault() time.Duration {
	return configutil.CoalesceDuration(c.ReadTimeout, DefaultReadTimeout)
}

// ReadHeaderTimeoutOrDefault gets a property.
func (c Config) ReadHeaderTimeoutOrDefault() time.Duration {
	return configutil.CoalesceDuration(c.ReadHeaderTimeout, DefaultReadHeaderTimeout)
}

// WriteTimeoutOrDefault gets a property.
func (c Config) WriteTimeoutOrDefault() time.Duration {
	return configutil.CoalesceDuration(c.WriteTimeout, DefaultWriteTimeout)
}

// IdleTimeoutOrDefault gets a property.
func (c Config) IdleTimeoutOrDefault() time.Duration {
	return configutil.CoalesceDuration(c.IdleTimeout, DefaultIdleTimeout)
}

// ShutdownGracePeriodOrDefault gets the shutdown grace period.
func (c Config) ShutdownGracePeriodOrDefault() time.Duration {
	return configutil.CoalesceDuration(c.ShutdownGracePeriod, DefaultShutdownGracePeriod)
}
