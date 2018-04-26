package main

import (
	"github.com/blend/go-sdk/logger"
	"github.com/blend/go-sdk/raft"
)

func main() {
	log := logger.All()

	cfg := raft.NewConfigFromEnv()
	r := raft.NewFromConfig(cfg).WithLogger(log)
	for _, remoteAddr := range cfg.GetPeers() {
		r = r.WithPeer(raft.NewRPCClient(remoteAddr).WithLogger(log))
	}

	if err := r.Start(); err != nil {
		log.SyncFatalExit(err)
	}

	select {}
}
