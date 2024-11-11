/*

Copyright (c) 2022 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package main

import (
	"time"

	"github.com/zpkg/blend-go-sdk/graceful"
	"github.com/zpkg/blend-go-sdk/logger"
	"github.com/zpkg/blend-go-sdk/web"
)

func main() {
	app := web.MustNew(web.OptLog(logger.All()))

	app.GET("/", func(_ *web.Ctx) web.Result {
		return web.NoContent
	}, web.WithTimeout(500*time.Millisecond), web.JSONProviderAsDefault)

	app.GET("/for/:duration", func(r *web.Ctx) web.Result {
		duration, err := web.DurationValue(r.RouteParam("duration"))
		if err != nil {
			return web.JSON.BadRequest(err)
		}
		time.Sleep(duration)
		return web.NoContent
	}, web.WithTimeout(5*time.Second), web.JSONProviderAsDefault)

	app.GET("/panic", func(_ *web.Ctx) web.Result {
		panic("ONLY A TEST")
	}, web.WithTimeout(500*time.Millisecond), web.JSONProviderAsDefault)

	if err := graceful.Shutdown(app); err != nil {
		logger.FatalExit(err)
	}
}
