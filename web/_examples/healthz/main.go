package main

import (
	"github.com/blend/go-sdk/env"
	"github.com/blend/go-sdk/logger"
	"github.com/blend/go-sdk/web"
)

func main() {
	log := logger.All().WithHeading("app")
	app := web.NewFromEnv().WithLogger(log)
	app.GET("/", func(r *web.Ctx) web.Result {
		return r.Text().Result("ok!")
	})

	hzLog := logger.All().WithHeading("healthz")
	hz := web.NewHealthz(app).WithLogger(hzLog)
	hzApp := web.New().WithLogger(hzLog).WithBindAddr(env.Env().String("HZ_BIND_ADDR", "127.0.0.1:8081")).WithHandler(hz)
	go func() {
		if err := web.StartWithGracefulShutdown(app); err != nil {
			logger.FatalExit(err)
		}
	}()
	if err := web.StartWithGracefulShutdown(hzApp); err != nil {
		logger.FatalExit(err)
	}
}
