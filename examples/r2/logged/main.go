/*

Copyright (c) 2022 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package main

import (
	"fmt"
	"os"

	"github.com/zpkg/blend-go-sdk/logger"
	"github.com/zpkg/blend-go-sdk/r2"
)

func main() {
	log := logger.MustNew(logger.OptAll())
	defer log.Close()

	_, err := r2.New("https://google.com/robots.txt",
		r2.OptHeaderValue("X-Sent-By", "go-sdk/request2"),
		r2.OptCookieValue("r2-ray-id", "example-stringdog01"),
		r2.OptLogResponse(log),
	).Discard()

	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
