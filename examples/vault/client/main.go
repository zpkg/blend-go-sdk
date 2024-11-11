/*

Copyright (c) 2022 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package main

import (
	"context"

	"github.com/zpkg/blend-go-sdk/logger"
	"github.com/zpkg/blend-go-sdk/vault"
)

func main() {
	log := logger.All()
	client, _ := vault.New(vault.OptConfigFromEnv(), vault.OptLog(log))

	key := "cubbyhole/willtest"

	ctx := context.Background()

	if err := client.Put(ctx, key, vault.Values{"value": "THE FOOOS"}); err != nil {
		log.Fatal(err)
		return
	}
	if err := client.Put(ctx, key, vault.Values{"value": "THE BUZZ"}); err != nil {
		log.Fatal(err)
		return
	}

	values, err := client.Get(ctx, key)
	if err != nil {
		log.Fatal(err)
		return
	}
	log.Infof("values: %#v", values)

	if err := client.Delete(ctx, key); err != nil {
		log.Fatal(err)
		return
	}
	log.Infof("~fin~")
}
