package main

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/blend/go-sdk/logger"
	"github.com/blend/go-sdk/web"
)

func main() {
	app := web.NewFromEnv().WithLogger(logger.All())

	app.ServeStaticCached("/cached", "_static")
	app.SetStaticMiddleware("/cached", web.SessionMiddleware(func(ctx *web.Ctx) web.Result {
		return ctx.Text().NotAuthorized()
	}, web.SessionReadLock))

	app.ServeStatic("/static", "_static")
	app.SetStaticMiddleware("/static", web.SessionMiddleware(func(ctx *web.Ctx) web.Result {
		return ctx.Text().NotAuthorized()
	}, web.SessionUnsafe))

	app.ServeStatic("/static_unauthed", "_static")

	app.Auth().SetValidateHandler(func(session *web.Session, state web.State) error {
		if session.UserID == "bailey" {
			return fmt.Errorf("bailey isn't allowed here")
		}
		return nil
	})

	app.Auth().SetLoginRedirectHandler(func(ctx *web.Ctx) *url.URL {
		u := *ctx.Request().URL
		u.Path = fmt.Sprintf("/login")
		return &u
	})

	app.Views().AddLiterals(`{{ define "login" }}<a href="/login/user_valid">Login Valid</a><br/><a href="/login/user_notvalid">Login Invalid</a>{{end}}`)
	app.GET("/login", func(r *web.Ctx) web.Result {
		return r.View().View("login", nil)
	})

	app.GET("/login/:userID", func(r *web.Ctx) web.Result {
		if r.Session() != nil {
			return r.RedirectWithMethodf("GET", "/")
		}

		userID := r.ParamString("userID")
		if !strings.HasSuffix(userID, "_valid") { //maximum security
			return r.Text().NotAuthorized()
		}
		_, err := r.Auth().Login(userID, r)
		if err != nil {
			return r.Text().InternalError(err)
		}
		return r.Text().Result("Logged In")
	}, web.SessionAware)

	app.GET("/logout", func(r *web.Ctx) web.Result {
		if r.Session() == nil {
			return r.Text().Result("Weren't logged in anyway.")
		}
		err := r.Auth().Logout(r)
		if err != nil {
			return r.Text().InternalError(err)
		}
		return r.Text().Result("Logged Out")
	}, web.SessionAware)

	app.GET("/", func(r *web.Ctx) web.Result {
		return r.Text().Result("Yep")
	}, web.SessionRequired)

	app.Logger().SyncFatalExit(app.Start())
}
