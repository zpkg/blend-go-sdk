package main

import (
	"github.com/blend/go-sdk/db"
	"github.com/blend/go-sdk/logger"
)

func main() {
	log := logger.NewFromEnv()
	conn := db.NewFromEnv().WithLogger(log)

	if err := conn.Open(); err != nil {
		log.SyncFatalExit(err)
	}

	log.SyncInfof("OK")
}
