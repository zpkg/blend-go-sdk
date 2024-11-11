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
	log := logger.Prod()
	app := web.MustNew(
		web.OptLog(log),
		web.OptConfigFromEnv(),
		web.OptUse(web.GZip),
	)
	app.GET("/", func(_ *web.Ctx) web.Result { return web.Text.Result("OK!") })
	if err := graceful.Shutdown(app); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
