package main

import (
	"net/http"

	"github.com/blend/go-sdk/logger"
	web "github.com/blendlabs/go-web"
)

func main() {
	app := web.New().WithLogger(logger.NewFromEnv())

	sf := web.NewCachedStaticFileServer(http.Dir("."))

	app.Static("/static/*filepath", "_static")
	app.StaticCached("/static_cached/*filepath", "_static")
	app.GET("/", func(r *web.Ctx) web.Result {
		return r.Static("index.html")
	})
	app.GET("/cached", func(r *web.Ctx) web.Result {
		cf, err := sf.GetCachedFile("index.html")
		if err != nil {
			return r.View().InternalError(err)
		}
		return cf
	})
	app.Start()
}
