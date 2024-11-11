/*

Copyright (c) 2022 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package main

import (
	"os"

	"github.com/zpkg/blend-go-sdk/graceful"
	"github.com/zpkg/blend-go-sdk/logger"
	"github.com/zpkg/blend-go-sdk/web"
)

func main() {
	app := web.MustNew(web.OptLog(logger.All()))
	app.Views.AddPaths(
		"_views/header.html",
		"_views/footer.html",
		"_views/index.html",
	)

	app.Views.FuncMap["foo"] = func() string {
		return "hello!"
	}

	if len(os.Getenv("LIVE_RELOAD")) > 0 {
		app.Views.LiveReload = true
	}

	app.GET("/", func(r *web.Ctx) web.Result {
		return r.Views.View("index", nil)
	})

	if err := graceful.Shutdown(app); err != nil {
		app.Log.Fatal(err)
		os.Exit(1)
	}
}
