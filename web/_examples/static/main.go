package main

import (
	"net/http"

	"github.com/blend/go-sdk/logger"
	"github.com/blend/go-sdk/web"
)

func main() {
	log := logger.MustNewFromEnv()
	app := web.New().WithLogger(log)
	csf := web.NewCachedStaticFileServer(http.Dir("."))

	app.ServeStatic("/static/*filepath", "_static")
	app.ServeStaticCached("/static_cached/*filepath", "_static")
	app.GET("/", func(r *web.Ctx) web.Result {
		return r.Static("index.html")
	})
	app.GET("/cached", func(r *web.Ctx) web.Result {
		return csf.ServeFile(r, "index.html")
	})
	log.SyncFatalExit(app.Start())
}
