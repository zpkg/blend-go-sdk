package main

import (
	"github.com/blend/go-sdk/logger"
	"github.com/blend/go-sdk/secrets"
)

func main() {
	log := logger.All()
	client := secrets.Must(secrets.NewFromConfig(secrets.NewConfigFromEnv())).WithLogger(log)

	if err := client.Put("secret/data/foo/bar", secrets.Values{"data": secrets.Values{"value": "THE FOOOS"}}); err != nil {
		log.SyncFatalExit(err)
	}
	if err := client.Put("secret/data/foo/baz", secrets.Values{"data": secrets.Values{"value": "THE BUZZ"}}); err != nil {
		log.SyncFatalExit(err)
	}

	value, err := client.Get("secret/data/foo/bar")
	if err != nil {
		log.SyncFatalExit(err)
	}
	log.Infof("value: %s", value)

	value, err = client.Get("secret/data/foo/bar")
	if err != nil {
		log.SyncFatalExit(err)
	}
	log.Infof("value: %s", value)

	if err := client.Delete("secret/data/foo/baz"); err != nil {
		log.SyncFatalExit(err)
	}
	if err := client.Delete("secret/data/foo/bar"); err != nil {
		log.SyncFatalExit(err)
	}
	log.Infof("~fin~")
}
