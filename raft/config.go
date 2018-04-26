package raft

import (
	"time"

	"github.com/blend/go-sdk/env"
	"github.com/blend/go-sdk/util"
	"github.com/blend/go-sdk/uuid"
)

const (
	// DefaultLeaderCheckTick is the tick rate for the leader check.
	DefaultLeaderCheckTick = 250 * time.Millisecond

	// DefaultHeartbeatTick is the tick rate for leaders to send heartbeats.
	DefaultHeartbeatTick = 250 * time.Millisecond

	// DefaultBindAddr is the default bind address.
	DefaultBindAddr = ":6060"
	// DefaultElectionTimeout is a default.
	DefaultElectionTimeout = 10 * time.Second
	// DefaultLeaderLeaseTimeout is a default.
	DefaultLeaderLeaseTimeout = 5 * time.Second

	// EnvVarIdentifier is an environment variable.
	EnvVarIdentifier = "RAFT_ID"
	// EnvVarBindAddr is an environment variable.
	EnvVarBindAddr = "RAFT_BIND_ADDR"
	// EnvVarPeers is an environment variable.
	EnvVarPeers = "RAFT_PEERS"
	// EnvVarElectionTimeout is an environment variable.
	EnvVarElectionTimeout = "RAFT_ELECTION_TIMEOUT"
	// EnvVarLeaderLeaseTimeout is an environment variable.
	EnvVarLeaderLeaseTimeout = "RAFT_LEADER_LEASE_TIMEOUT"
)

// NewConfigFromEnv creates a new config from environment variables.
func NewConfigFromEnv() *Config {
	var cfg Config
	if err := env.Env().ReadInto(&cfg); err != nil {
		panic(err)
	}
	return &cfg
}

// Config is a node config.
type Config struct {
	Identifier         string        `json:"identifier,omitempty" yaml:"identifier,omitempty" env:"RAFT_ID"`
	BindAddr           string        `json:"bindAddr,omitempty" yaml:"bindAddr,omitempty" env:"RAFT_BIND_ADDR"`
	Peers              []string      `json:"peers,omitempty" yaml:"peers,omitempty" env:"RAFT_PEERS,csv"`
	ElectionTimeout    time.Duration `json:"electionTimeout,omitempty" yaml:"electionTimeout,omitempty" env:"RAFT_ELECTION_TIMEOUT"`
	LeaderLeaseTimeout time.Duration `json:"leaderLeaseTimeout,omitempty" yaml:"leaderLeaseTimeout,omitempty" env:"RAFT_LEADER_LEASE_TIMEOUT"`
}

// GetIdentifier gets a field or a default.
func (c Config) GetIdentifier(inherited ...string) string {
	return util.Coalesce.String(c.Identifier, uuid.V4().String(), inherited...)
}

// GetBindAddr gets a field or a default.
func (c Config) GetBindAddr(inherited ...string) string {
	return util.Coalesce.String(c.BindAddr, DefaultBindAddr, inherited...)
}

// GetPeers gets a field or a default.
func (c Config) GetPeers(inherited ...[]string) []string {
	return util.Coalesce.Strings(c.Peers, nil, inherited...)
}

// GetElectionTimeout gets a field or a default.
func (c Config) GetElectionTimeout(inherited ...time.Duration) time.Duration {
	return util.Coalesce.Duration(c.ElectionTimeout, DefaultElectionTimeout, inherited...)
}

// GetLeaderLeaseTimeout gets a field or a default.
func (c Config) GetLeaderLeaseTimeout(inherited ...time.Duration) time.Duration {
	return util.Coalesce.Duration(c.LeaderLeaseTimeout, DefaultLeaderLeaseTimeout, inherited...)
}
