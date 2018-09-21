package main

import (
	"time"

	"github.com/blend/go-sdk/env"
	"github.com/blend/go-sdk/logger"
	"github.com/blend/go-sdk/web"
)

func main() {
	log := logger.All().WithHeading("app")
	app := web.MustNewFromEnv().WithLogger(log)

	app.GET("/", func(r *web.Ctx) web.Result {
		return r.Text().Result("ok!")
	})

	// create a healthz and host our app within it.
	hz := web.NewHealthz(app).
		WithBindAddr(env.Env().String("HZ_BIND_ADDR", "127.0.0.1:8081")).
		WithGracePeriod(30 * time.Second).
		WithFailureThreshold(3).
		WithLogger(logger.All().WithHeading("healthz"))

	// start the hz and the child app, ideally they have separate bind addrs.
	if err := web.StartWithGracefulShutdown(hz); err != nil {
		logger.FatalExit(err)
	}
}
