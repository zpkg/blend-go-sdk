/*

Copyright (c) 2022 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package main

import (
	"time"

	"github.com/zpkg/blend-go-sdk/logger"
)

func main() {
	log := logger.MustNew(logger.OptConfigFromEnv())
	tick := time.Tick(time.Second)
	for range tick {
		log.Infof("it's %s", time.Now().Format(time.RFC3339))
	}
}
