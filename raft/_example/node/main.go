package main

import (
	"flag"
	"strings"
	"time"

	"github.com/blend/go-sdk/logger"
	"github.com/blend/go-sdk/raft"
)

var (
	id       = flag.String("id", "", "The raft worker id")
	bindAddr = flag.String("bind-addr", "", "Bind address")
)

func main() {
	flag.Parse()
	log := logger.All()

	cfg := raft.NewConfigFromEnv().WithID(*id).WithBindAddr(*bindAddr)

	r := raft.NewFromConfig(cfg).WithLogger(log)
	r.WithServer(raft.NewRPCServer().WithBindAddr(cfg.GetBindAddr()))
	log.WithHeading(r.ID())
	for _, remoteAddr := range cfg.GetPeers() {
		// don't add the peer if it's outself.
		if !strings.HasSuffix(strings.TrimSpace(remoteAddr), cfg.GetBindAddr()) {
			log.SyncDebugf("adding peer %s", remoteAddr)
			r = r.WithPeer(raft.NewRPCClient(remoteAddr))
		} else {
			log.SyncDebugf("skipping %s", remoteAddr)
		}
	}

	// wait until the next second

	now := time.Now().UTC()
	next := now.Add(time.Second)
	next = next.Add(-time.Duration(next.Nanosecond()))

	log.Debugf("synchronized start; sleeping for %v", next.Sub(now))
	time.Sleep(next.Sub(now))

	if err := r.Start(); err != nil {
		log.SyncFatalExit(err)
	}

	select {}
}
