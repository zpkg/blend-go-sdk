package main

import (
	"github.com/blend/go-sdk/env"
	"github.com/blend/go-sdk/logger"
	"github.com/blend/go-sdk/web"
)

/*
func main() {
	log := logger.All().WithHeading("app")
	app := web.NewFromEnv().WithLogger(log)
	app.GET("/", func(r *web.Ctx) web.Result {
		return r.Text().Result("ok!")
	})

	hzLog := logger.All().WithHeading("healthz")


	hz := web.NewHealthz(app).WithLogger(hzLog)
	hzApp := web.New().
		WithLogger(hzLog).
		WithBindAddr(env.Env().String("HZ_BIND_ADDR", "127.0.0.1:8081")).
		WithHandler(hz).
		WithChildApp(hz.App()).
		WithPreShutdownCallBack(func(_ *web.App) error {
			hz.SetReady(false)           // set fail on probes
			time.Sleep(10 * time.Second) // wait for probes to fail
			return nil
		})

	if err := web.StartWithGracefulShutdown(hzApp); err != nil {
		logger.FatalExit(err)
	}
}
*/

func main() {
	log := logger.All()
	app := web.NewFromEnv().WithLogger(log)

	app.GET("/", func(r *web.Ctx) web.Result {
		return r.Text().Result("ok!")
	})
	app.GET("/shutdown", func(r *web.Ctx) web.Result {
		app.Shutdown()
		return nil
	})

	// create a healthz and host our app within it.
	hz := web.NewHealthz(app).WithBindAddr(env.Env().String("HZ_BIND_ADDR", "127.0.0.1:8081"))

	// start the hz and the child app, ideally they have separate bind addrs.
	if err := web.StartWithGracefulShutdown(hz); err != nil {
		logger.FatalExit(err)
	}
}
