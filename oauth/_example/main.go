package main

import (
	"fmt"
	"net/url"

	"github.com/blend/go-sdk/configutil"
	"github.com/blend/go-sdk/logger"
	google "github.com/blend/go-sdk/oauth"
	"github.com/blend/go-sdk/util"
	"github.com/blend/go-sdk/web"
)

// Config is the app config.
type Config struct {
	Logger     logger.Config `json:"logger" yaml:"logger"`
	GoogleAuth google.Config `json:"googleAuth" yaml:"googleAuth"`
	Web        web.Config    `json:"web" yaml:"web"`
}

func main() {
	var cfg Config
	if err := configutil.Read(&cfg, "./config.yml"); err != nil {
		logger.All().SyncFatalExit(err)
	}

	log := logger.NewFromConfig(&cfg.Logger)
	app := web.NewFromConfig(&cfg.Web).WithLogger(log)
	if app.Err() != nil {
		log.SyncFatalExit(app.Err())
	}

	// create the oauth manager from the section in the config.
	// also, give it a secret
	oauth := google.NewFromConfig(&cfg.GoogleAuth).WithSecret(util.Crypto.MustCreateKey(32))

	// if we haven't set a redirect uri explicitly in the config
	// add one based on the base url for the app and the controller route.
	if len(oauth.RedirectURI()) == 0 {
		if app.BaseURL() == nil {
			log.SyncFatalExit(fmt.Errorf("web:baseURL must be set in the config"))
		}
		oauth = oauth.WithRedirectURI(fmt.Sprintf("%s/oauth/google", app.BaseURL().String()))
	}

	// check if there are issues with the oauth config.
	// we do this here in case we set the redirectURI from the web app base url.
	if err := oauth.ValidateConfig(); err != nil {
		log.SyncFatalExit(err)
	}

	// if a route is marked as `web.SessionRequired`, how should the auth manager
	// redirect the user? typically this should punt the user to the login page.
	// this should also handle remembering where the user was trying to go originally.
	app.Auth().WithLoginRedirectHandler(func(ctx *web.Ctx) *url.URL {
		original := ctx.Request().URL
		newURL := *original
		newURL.Path = "/login"

		query := url.Values{}
		query.Add("redirect", original.String())
		newURL.ForceQuery = true
		newURL.RawQuery = query.Encode()
		return &newURL
	})

	app.Views().AddLiterals(`{{ define "login" }}<a href={{ .ViewModel.OAuthURL }}>Login</a>{{end}}`)

	app.GET("/", func(r *web.Ctx) web.Result {
		return r.JSON().Result("ok!")
	}, web.SessionRequired)

	app.GET("/login", func(r *web.Ctx) web.Result {
		// corner case; we're already logged in
		if r.Session() != nil {
			r.RedirectWithMethodf("GET", "/")
		}

		authURL, err := oauth.OAuthURL(r.ParamString("redirect"))
		if err != nil {
			return r.Text().InternalError(err)
		}

		return r.View().View("login", map[string]interface{}{
			"OAuthURL": authURL,
		})
	}, web.SessionAware)

	app.GET("/oauth/google", func(r *web.Ctx) web.Result {
		// corner case; we're already logged in
		if r.Session() != nil {
			return r.RedirectWithMethodf("GET", "/")
		}

		// finish the oauth process
		// this checks the state nonce and
		result, err := oauth.Finish(r.Request())
		if err != nil {
			return r.Text().InternalError(err)
		}

		_, err = r.Auth().Login(result.Profile.Email, r)
		if err != nil {
			return r.Text().InternalError(err)
		}

		return r.RedirectWithMethodf("GET", "/")
	}, web.SessionAware)

	app.GET("/logout", func(r *web.Ctx) web.Result {
		err := r.Auth().Logout(r)
		if err != nil {
			return r.Text().InternalError(err)
		}

		return r.RedirectWithMethodf("GET", "/")
	}, web.SessionRequired)

	log.SyncFatalExit(app.Start())
}
