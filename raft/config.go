package raft

import (
	"time"

	"github.com/blend/go-sdk/env"
	"github.com/blend/go-sdk/util"
	"github.com/blend/go-sdk/uuid"
)

const (
	// DefaultBindAddr is the default bind address.
	DefaultBindAddr = ":6060"
	// DefaultHeartbeatTimeout is a default.
	DefaultHeartbeatTimeout = time.Second
	// DefaultElectionTimeout is a default.
	DefaultElectionTimeout = time.Second
	// DefaultLeaderLeaseTimeout is a default.
	DefaultLeaderLeaseTimeout = 5 * time.Second
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
	Identifier         string        `yaml:"identifier" env:"RAFT_ID"`
	BindAddr           string        `yaml:"bindAddr" env:"RAFT_BIND_ADDR"`
	Peers              []string      `yaml:"peers" env:"RAFT_PEERS,csv"`
	HeartbeatTimeout   time.Duration `yaml:"heatbeatTimeout" env:"RAFT_HEARTBEAT_TIMEOUT"`
	ElectionTimeout    time.Duration `yaml:"electionTimeout" env:"RAFT_ELECTION_TIMEOUT"`
	LeaderLeaseTimeout time.Duration `yaml:"leaderLeaseTimeout" env:"RAFT_LEADER_LEASE_TIMEOUT"`
	StartAsLeader      *bool         `yaml:"startAsLeader" env:"RAFT_START_LEADER"`
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

// GetHeartbeatTimeout gets a field or a default.
func (c Config) GetHeartbeatTimeout(inherited ...time.Duration) time.Duration {
	return util.Coalesce.Duration(c.HeartbeatTimeout, DefaultHeartbeatTimeout, inherited...)
}

// GetElectionTimeout gets a field or a default.
func (c Config) GetElectionTimeout(inherited ...time.Duration) time.Duration {
	return util.Coalesce.Duration(c.ElectionTimeout, DefaultElectionTimeout, inherited...)
}

// GetLeaderLeaseTimeout gets a field or a default.
func (c Config) GetLeaderLeaseTimeout(inherited ...time.Duration) time.Duration {
	return util.Coalesce.Duration(c.LeaderLeaseTimeout, DefaultLeaderLeaseTimeout, inherited...)
}

// GetStartAsLeader gets a field or a default.
func (c Config) GetStartAsLeader(inherited ...bool) bool {
	return util.Coalesce.Bool(c.StartAsLeader, false, inherited...)
}
