package main

import (
	"flag"

	"github.com/blend/go-sdk/logger"
	"github.com/blend/go-sdk/raft"
)

var (
	id       = flag.String("id", "", "The raft worker id")
	selfAddr = flag.String("self-addr", "", "Our address in the peer list")
	bindAddr = flag.String("bind-addr", "", "Bind address")
)

func main() {
	flag.Parse()
	log := logger.All()

	cfg := raft.NewConfigFromEnv().WithID(*id).WithSelfAddr(*selfAddr).WithBindAddr(*bindAddr)

	r := raft.NewFromConfig(cfg).WithLogger(log)
	log.WithHeading(r.ID())
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
