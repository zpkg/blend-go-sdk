package main

import (
	"github.com/blend/go-sdk/logger"
	"github.com/blend/go-sdk/raft"
)

func main() {
	log := logger.All()

	r := raft.NewFromConfig(raft.NewConfigFromEnv()).WithLogger(log)

	for _, remoteAddr := range r.Config().GetPeers() {
		r = r.WithPeer(raft.NewClient(remoteAddr).WithLogger(log))
	}

	if err := r.Start(); err != nil {
		log.SyncFatalExit(err)
	}

	select {}
}
