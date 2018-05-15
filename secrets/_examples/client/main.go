package main

import (
	"github.com/blend/go-sdk/logger"
	"github.com/blend/go-sdk/secrets"
)

func main() {
	log := logger.All()
	client := secrets.Must(secrets.NewFromEnv()).WithLogger(log)

	if err := client.Put("foo/bar", secrets.Values{"value": "THE FOOOS"}); err != nil {
		log.SyncFatalExit(err)
	}
	if err := client.Put("foo/baz", secrets.Values{"value": "THE BUZZ"}); err != nil {
		log.SyncFatalExit(err)
	}

	values, err := client.Get("foo/bar")
	if err != nil {
		log.SyncFatalExit(err)
	}
	log.Infof("values: %#v", values)

	values, err = client.Get("foo/baz")
	if err != nil {
		log.SyncFatalExit(err)
	}
	log.Infof("values: %#v", values)

	if err := client.Delete("foo/baz"); err != nil {
		log.SyncFatalExit(err)
	}
	if err := client.Delete("foo/bar"); err != nil {
		log.SyncFatalExit(err)
	}
	log.Infof("~fin~")
}
