package main

import (
	"github.com/blend/go-sdk/env"
	"github.com/blend/go-sdk/logger"
	"github.com/blend/go-sdk/web"
)

func main() {
	log := logger.All()
	app := web.NewFromEnv().WithLogger(log)

	app.GET("/", func(r *web.Ctx) web.Result {
		return r.Text().Result("ok!")
	})

	// create a healthz and host our app within it.
	hz := web.NewHealthz(app).
		WithBindAddr(env.Env().String("HZ_BIND_ADDR", "127.0.0.1:8081")).
		WithGracePeriodSeconds(30).
		WithLogger(logger.All().WithHeading("healthz"))

	// start the hz and the child app, ideally they have separate bind addrs.
	if err := web.StartWithGracefulShutdown(hz); err != nil {
		logger.FatalExit(err)
	}
}
