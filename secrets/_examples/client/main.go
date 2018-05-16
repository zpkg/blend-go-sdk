package main

import (
	"github.com/blend/go-sdk/logger"
	"github.com/blend/go-sdk/secrets"
)

func main() {
	log := logger.All()
	client := secrets.Must(secrets.NewFromEnv()).WithLogger(log)

	key := "cubbyhole/willtest"

	if err := client.Put(key, secrets.Values{"value": "THE FOOOS"}); err != nil {
		log.SyncFatalExit(err)
	}
	if err := client.Put(key, secrets.Values{"value": "THE BUZZ"}); err != nil {
		log.SyncFatalExit(err)
	}

	values, err := client.Get(key)
	if err != nil {
		log.SyncFatalExit(err)
	}
	log.Infof("values: %#v", values)

	if err := client.Delete(key); err != nil {
		log.SyncFatalExit(err)
	}
	log.Infof("~fin~")
}
