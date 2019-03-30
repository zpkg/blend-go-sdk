package web

import (
	"crypto/tls"
	"fmt"
	"net/http"

	"github.com/blend/go-sdk/logger"
)

// AppOption is an option for an app.
type AppOption func(*App)

// OptConfig sets the config.
func OptConfig(cfg *Config) AppOption {
	return func(a *App) {
		a.Config = cfg
		a.Auth = NewAuthManager(cfg)
		a.Views = NewViewCache(OptViewCacheConfig(&cfg.Views))
	}
}

// OptBindAddr sets the config bind address
func OptBindAddr(bindAddr string) AppOption {
	return func(a *App) {
		a.Config.BindAddr = bindAddr
	}
}

// OptPort sets the config bind address
func OptPort(port int32) AppOption {
	return func(a *App) {
		a.Config.Port = port
		a.Config.BindAddr = fmt.Sprintf(":%v", port)
	}
}

// OptLog sets the logger.
func OptLog(log logger.Log) AppOption {
	return func(a *App) { a.Log = log }
}

// OptServer sets the underlying server.
func OptServer(server *http.Server) AppOption {
	return func(a *App) { a.Server = server }
}

// OptAuth sets the auth manager.
func OptAuth(auth *AuthManager) AppOption {
	return func(a *App) { a.Auth = auth }
}

// OptViews sets the view cache.
func OptViews(views *ViewCache) AppOption {
	return func(a *App) { a.Views = views }
}

// OptHandler sets the underlying handler
func OptHandler(handler http.Handler) AppOption {
	return func(a *App) { a.Handler = handler }
}

// OptTLSConfig sets the tls config.
func OptTLSConfig(cfg *tls.Config) AppOption {
	return func(a *App) { a.TLSConfig = cfg }
}

// OptDefaultHeader sets a default header.
func OptDefaultHeader(key, value string) AppOption {
	return func(a *App) {
		if a.DefaultHeaders == nil {
			a.DefaultHeaders = make(map[string]string)
		}
		a.DefaultHeaders[key] = value
	}
}

// OptDefaultHeaders sets default headers.
func OptDefaultHeaders(headers map[string]string) AppOption {
	return func(a *App) { a.DefaultHeaders = headers }
}

// OptUse adds to the default middleware.
func OptUse(m Middleware) AppOption {
	return func(a *App) { a.DefaultMiddleware = append(a.DefaultMiddleware, m) }
}

// OptNotFoundHandler sets default headers.
func OptNotFoundHandler(action Action) AppOption {
	return func(a *App) { a.NotFoundHandler = a.RenderAction(action) }
}
