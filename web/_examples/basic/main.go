package main

import (
	"fmt"

	"github.com/blend/go-sdk/logger"
	"github.com/blend/go-sdk/web"
)

func main() {
	app := web.MustNewFromEnv().WithBindAddr(":8080").WithLogger(logger.All())
	app.GET("/", func(r *web.Ctx) web.Result {
		return r.Text().Result("ok!")
	})

	app.POST("/reparse", func(r *web.Ctx) web.Result {
		body, err := r.PostBody()
		if err != nil {
			return r.Text().BadRequest(err)
		}
		if len(body) == 0 {
			return r.Text().BadRequest(fmt.Errorf("empty body"))
		}
		return r.Text().Result(web.StringValue(r.Param("foo")))
	})
	logger.All().SyncFatalExit(app.Start())
}
