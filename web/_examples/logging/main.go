package main

import (
	"fmt"

	"github.com/blend/go-sdk/logger"
	"github.com/blend/go-sdk/web"
)

func main() {
	log := logger.MustNewFromEnv()
	app := web.MustNewFromEnv().WithLogger(log)
	app.GET("/", func(r *web.Ctx) web.Result {
		return r.Text().Result("foo")
	})

	log.Listen(logger.HTTPRequest, logger.DefaultListenerName, logger.NewHTTPRequestEventListener(func(wre *logger.HTTPRequestEvent) {
		fmt.Printf("Route: %s\n", wre.Route())
	}))

	log.SyncFatalExit(app.Start())
}
