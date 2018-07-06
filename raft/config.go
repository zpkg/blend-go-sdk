package raft

import (
	"time"

	"github.com/blend/go-sdk/env"
	"github.com/blend/go-sdk/util"
)

const (
	// DefaultLeaderCheckInterval is the tick rate for the leader check.
	DefaultLeaderCheckInterval = 2000 * time.Millisecond
	// DefaultHeartbeatInterval is the tick rate for leaders to send heartbeats.
	DefaultHeartbeatInterval = 1000 * time.Millisecond
	// DefaultElectionTimeout is a default.
	DefaultElectionTimeout = 5 * DefaultHeartbeatInterval
	// DefaultElectionBackoffTimeout is a default.
	DefaultElectionBackoffTimeout = DefaultElectionTimeout
	// DefaultPeerDialTimeout is the default peer dial timeout.
	DefaultPeerDialTimeout = time.Second

	// DefaultBindAddr is the default rpc server bind address.
	DefaultBindAddr = ":6060"

	// EnvVarIdentifier is an environment variable.
	EnvVarIdentifier = "RAFT_ID"
	// EnvVarBindAddr is an environment variable.
	EnvVarBindAddr = "RAFT_BIND_ADDR"
	// EnvVarPeers is an environment variable.
	EnvVarPeers = "RAFT_PEERS"
	// EnvVarElectionTimeout is an environment variable.
	EnvVarElectionTimeout = "RAFT_ELECTION_TIMEOUT"
	// EnvVarRaftPeerDialTimeout is an environment variable.
	EnvVarRaftPeerDialTimeout = "RAFT_PEER_DIAL_TIMEOUT"
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
	ID                  string        `json:"id,omitempty" yaml:"id,omitempty" env:"RAFT_ID"`
	BindAddr            string        `json:"bindAddr,omitempty" yaml:"bindAddr,omitempty" env:"RAFT_BIND_ADDR"`
	Peers               []string      `json:"peers,omitempty" yaml:"peers,omitempty" env:"RAFT_PEERS,csv"`
	HeartbeatInterval   time.Duration `json:"heartbeatInterval,omitempty" yaml:"heartbeatInterval,omitempty" env:"RAFT_HEARTBEAT_INTERVAL"`
	LeaderCheckInterval time.Duration `json:"leaderCheckInterval,omitempty" yaml:"leaderCheckInterval,omitempty" env:"RAFT_LEADER_CHECK_INTERVAL"`
	ElectionTimeout     time.Duration `json:"electionTimeout,omitempty" yaml:"electionTimeout,omitempty" env:"RAFT_ELECTION_TIMEOUT"`
	PeerDialTimeout     time.Duration `json:"peerDialTimeout,omitempty" yaml:"peerDialTimeout,omitempty" env:"RAFT_PEER_DIAL_TIMEOUT"`
}

// GetID gets a field or a default.
func (c Config) GetID(inherited ...string) string {
	return util.Coalesce.String(c.ID, "", inherited...)
}

// WithID sets the a property if the value is set.
func (c *Config) WithID(value string) *Config {
	if len(value) > 0 {
		c.ID = value
	}
	return c
}

// GetBindAddr gets a field or a default.
func (c Config) GetBindAddr(inherited ...string) string {
	return util.Coalesce.String(c.BindAddr, DefaultBindAddr, inherited...)
}

// WithBindAddr sets the a property if the value is set.
func (c *Config) WithBindAddr(value string) *Config {
	if len(value) > 0 {
		c.BindAddr = value
	}
	return c
}

// GetPeers gets a field or a default.
func (c Config) GetPeers(inherited ...[]string) []string {
	return util.Coalesce.Strings(c.Peers, nil, inherited...)
}

// GetHeartbeatInterval gets a field or a default.
func (c Config) GetHeartbeatInterval(inherited ...time.Duration) time.Duration {
	return util.Coalesce.Duration(c.HeartbeatInterval, DefaultHeartbeatInterval, inherited...)
}

// GetLeaderCheckInterval gets a field or a default.
func (c Config) GetLeaderCheckInterval(inherited ...time.Duration) time.Duration {
	return util.Coalesce.Duration(c.LeaderCheckInterval, DefaultLeaderCheckInterval, inherited...)
}

// GetElectionTimeout gets a field or a default.
func (c Config) GetElectionTimeout(inherited ...time.Duration) time.Duration {
	return util.Coalesce.Duration(c.ElectionTimeout, DefaultElectionTimeout, inherited...)
}

// GetPeerDialTimeout gets a field or a default.
func (c Config) GetPeerDialTimeout(inherited ...time.Duration) time.Duration {
	return util.Coalesce.Duration(c.PeerDialTimeout, DefaultPeerDialTimeout, inherited...)
}
