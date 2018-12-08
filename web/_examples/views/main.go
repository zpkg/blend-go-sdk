package main

import (
	"os"

	"github.com/blend/go-sdk/logger"
	"github.com/blend/go-sdk/web"
)

func main() {
	log := logger.MustNewFromEnv()
	app := web.MustNewFromEnv().WithLogger(log)
	app.Views().AddPaths(
		"_views/header.html",
		"_views/footer.html",
		"_views/index.html",
	)

	app.Views().FuncMap()["foo"] = func() string {
		return "hello!"
	}

	if len(os.Getenv("LIVE_RELOAD")) > 0 {
		app.Views().WithCached(false)
	}

	app.GET("/", func(r *web.Ctx) web.Result {
		return r.View().View("index", nil)
	})
	if err := web.GracefulShutdown(app); err != nil {
		log.SyncFatalExit(err)
	}
}
