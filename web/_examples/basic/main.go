package main

import (
	"github.com/blend/go-sdk/logger"
	"github.com/blend/go-sdk/web"
)

func main() {
	app := web.NewFromEnv().WithBindAddr(":8080")
	app.GET("/", func(r *web.Ctx) web.Result {
		return r.Text().Result("ok!")
	})
	logger.All().SyncFatalExit(app.Start())
}
