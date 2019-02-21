package main

import (
	"context"

	"github.com/blend/go-sdk/logger"
	"github.com/blend/go-sdk/secrets"
)

func main() {
	log := logger.Sync()
	client := secrets.Must(secrets.NewFromEnv()).WithLogger(log)

	key := "cubbyhole/willtest"

	ctx := context.Background()

	if err := client.Put(ctx, key, secrets.Values{"value": "THE FOOOS"}); err != nil {
		log.SyncFatalExit(err)
	}
	if err := client.Put(ctx, key, secrets.Values{"value": "THE BUZZ"}); err != nil {
		log.SyncFatalExit(err)
	}

	values, err := client.Get(ctx, key)
	if err != nil {
		log.SyncFatalExit(err)
	}
	log.SyncInfof("values: %#v", values)

	if err := client.Delete(ctx, key); err != nil {
		log.SyncFatalExit(err)
	}
	log.SyncInfof("~fin~")
}
