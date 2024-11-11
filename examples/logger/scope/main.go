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
	all := logger.MustNew(logger.OptAll())
	go func(log logger.Log) {
		ticker := time.Tick(500 * time.Millisecond)
		for {
			<-ticker
			log.Infof("this is foo")
		}
	}(all.WithPath("foo"))

	go func(log logger.Log) {
		ticker := time.Tick(500 * time.Millisecond)
		for {
			<-ticker
			log.Infof("this is bar")
		}
	}(all.WithPath("bar"))

	select {}
}
