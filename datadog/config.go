package datadog

import (
	"fmt"

	"github.com/blend/go-sdk/configutil"
	"github.com/blend/go-sdk/env"
)

const (
	// DefaultDatadogBufferDepth is the default number of statsd messages to buffer.
	DefaultDatadogBufferDepth = 128
)

// MustNewConfigFromEnv creates a new config from the environment and panics on error.
func MustNewConfigFromEnv() (config *Config) {
	var err error
	if config, err = NewConfigFromEnv(); err != nil {
		panic(err)
	}
	return
}

// NewConfigFromEnv returns a new config from the env.
func NewConfigFromEnv() (*Config, error) {
	var config Config
	if err := env.Env().ReadInto(&config); err != nil {
		return nil, err
	}
	return &config, nil
}

// Config is the datadog config.
type Config struct {
	// Hostname is the dns name or ip of the datadog collector.
	Hostname string `json:"hostname,omitempty" yaml:"hostname,omitempty" env:"DATADOG_HOSTNAME"`
	// Port is the port of the datadog collector.
	Port string `json:"port,omitempty" yaml:"port,omitempty" env:"DATADOG_PORT"`
	// TracePort is the port of the datadog apm collector.
	TracePort string `json:"tracePort,omitempty" yaml:"tracePort,omitempty" env:"DATADOG_TRACE_PORT"`
	// TracingEnabled returns if we should use tracing or not.
	TracingEnabled *bool `json:"tracingEnabled" yaml:"tracingEnabled" env:"DATADOG_APM_ENABLED"`
	// Buffered indicates if we should buffer statsd messages or not.
	Buffered *bool `json:"buffered,omitempty" yaml:"buffered,omitempty" env:"DATADOG_BUFFERED"`
	// BufferDepth is the depth of the buffer for datadog events.
	// A zero value implies an unbuffered client.
	BufferDepth int `json:"bufferDepth,omitempty" yaml:"bufferDepth,omitempty" env:"DATADOG_BUFFER_DEPTH"`
	// Namespace is an optional namespace.
	Namespace string `json:"namespace,omitempty" yaml:"namespace,omitempty" env:"DATADOG_NAMESPACE"`
	// DefaultTags are the default tags associated with any stat metric.
	DefaultTags []string `json:"defaultTags,omitempty" yaml:"defaultTags,omitempty" env:"DATADOG_TAGS,csv"`
}

// IsZero returns if the config is unset.
func (c Config) IsZero() bool {
	return len(c.HostnameOrDefault()) == 0
}

// HostnameOrDefault returns the datadog hostname.
func (c Config) HostnameOrDefault(defaults ...string) string {
	return configutil.CoalesceString(c.Hostname, "", defaults...)
}

// PortOrDefault returns the datadog port.
func (c Config) PortOrDefault(defaults ...string) string {
	return configutil.CoalesceString(c.Port, DefaultPort, defaults...)
}

// TracePortOrDefault returns the datadog trace port.
func (c Config) TracePortOrDefault(defaults ...string) string {
	return configutil.CoalesceString(c.TracePort, DefaultTracePort, defaults...)
}

// TracingEnabledOrDefault returns if tracing is enabled.
func (c Config) TracingEnabledOrDefault() bool {
	return configutil.CoalesceBool(c.TracingEnabled, DefaultTracingEnabled)
}

// Host returns the datadog collector host:port string.
func (c Config) Host() string {
	return fmt.Sprintf("%s:%s", c.HostnameOrDefault(), c.PortOrDefault())
}

// TraceHostOrDefault returns the datadog trace collector host:port string.
func (c Config) TraceHostOrDefault() string {
	return fmt.Sprintf("%s:%s", c.HostnameOrDefault(), c.TracePortOrDefault())
}

// BufferedOrDefault returns if the client should buffer messages or not.
func (c Config) BufferedOrDefault(defaults ...bool) bool {
	return configutil.CoalesceBool(c.Buffered, false, defaults...)
}

// BufferDepthOrDefault returns the buffer depth.
func (c Config) BufferDepthOrDefault(defaults ...int) int {
	return configutil.CoalesceInt(c.BufferDepth, DefaultDatadogBufferDepth, defaults...)
}

// NamespaceOrDefault returns the default prefix for metric names.
func (c Config) NamespaceOrDefault(defaults ...string) string {
	return configutil.CoalesceString(c.Namespace, "", defaults...)
}

// DefaultTagsOrDefault returns default tags for the client.
func (c Config) DefaultTagsOrDefault(defaults ...[]string) []string {
	return configutil.CoalesceStrings(c.DefaultTags, nil, defaults...)
}
