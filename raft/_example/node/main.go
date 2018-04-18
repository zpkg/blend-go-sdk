package main

import (
	"github.com/blend/go-sdk/logger"
	"github.com/blend/go-sdk/raft"
)

func main() {
	log := logger.New()

	r := raft.NewFromConfig(raft.NewConfigFromEnv())
}
