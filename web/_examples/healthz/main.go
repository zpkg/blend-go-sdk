package main

import (
	"github.com/blend/go-sdk/logger"
	web "github.com/blendlabs/go-web"
)

func main() {
	log := logger.All().WithLabel("app")
	app := web.NewFromEnv().WithLogger(log)
	app.GET("/", func(r *web.Ctx) web.Result {
		return r.Text().Result("ok!")
	})
	web.HealthzHost(app, web.NewHealthzFromEnv(app).WithLogger(log))
}
