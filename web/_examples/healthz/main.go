package main

import (
	"github.com/blend/go-sdk/logger"
	"github.com/blend/go-sdk/web"
)

func main() {
	log := logger.All().WithHeading("app")
	app := web.NewFromEnv().WithLogger(log)
	app.GET("/", func(r *web.Ctx) web.Result {
		return r.Text().Result("ok!")
	})

	hz := web.NewHealthzFromEnv(app).WithLogger(log)
	go func() {
		if err := web.StartWithGracefulShutdown(app); err != nil {
			logger.FatalExit(err)
		}
	}()
	if err := web.StartWithGracefulShutdown(web.New().WithServer(hz.Server())); err != nil {
		logger.FatalExit(err)
	}
}
