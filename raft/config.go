package raft

import "time"

// Config is a node config.
type Config struct {
	Identifier         string        `yaml:"identifier" env:"RAFT_IDENTIFIER"`
	Peers              []string      `yaml:"peers" env:"RAFT_PEERS"`
	HeartbeatTimeout   time.Duration `yaml:"heatbeatTimeout" env:"RAFT_HEARTBEAT_TIMEOUT"`
	ElectionTimeout    time.Duration `yaml:"electionTimeout" env:"RAFT_ELECTION_TIMEOUT"`
	LeaderLeaseTimeout time.Duration `yaml:"leaderLeaseTimeout" env:"RAFT_LEADER_LEASE_TIMEOUT"`
	StartAsLeader      *bool         `yaml:"startAsLeader" env:"RAFT_START_LEADER"`
}
