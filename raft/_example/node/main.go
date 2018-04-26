package main

import (
	"github.com/blend/go-sdk/logger"
	"github.com/blend/go-sdk/raft"
)

func main() {
	log := logger.All().WithDisabled(logger.Error)

	cfg := raft.NewConfigFromEnv()
	r := raft.NewFromConfig(cfg).WithLogger(log)
	for _, remoteAddr := range cfg.GetPeers() {
		if !r.IsSelf(remoteAddr) {
			log.SyncDebugf("adding peer %s", remoteAddr)
			r = r.WithPeer(raft.NewRPCClient(remoteAddr).WithLogger(log))
		} else {
			log.SyncDebugf("skipping %s", remoteAddr)
		}
	}

	if err := r.Start(); err != nil {
		log.SyncFatalExit(err)
	}

	select {}
}
