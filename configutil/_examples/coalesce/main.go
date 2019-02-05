package main

import (
	"flag"
	"fmt"

	cfg "github.com/blend/go-sdk/configutil"
	"github.com/blend/go-sdk/logger"
)

var (
	flagTarget      = flag.String("target", "", "The target URL")
	flagEnvironment = flag.String("env", "", "The current environment")
)

// Config is a sample config.
type Config struct {
	Target      string `yaml:"target"`
	Environment string `yaml:"env"`
}

// Resolve parses the config and sets values from predefined sources.
func (c *Config) Resolve() error {
	return cfg.AnyError(
		cfg.Set(&c.Target, cfg.Const(*flagTarget), cfg.Env("TARGET"), cfg.Const(c.Target), cfg.Const("https://google.com/robots.txt")),
		cfg.Set(&c.Environment, cfg.Const(*flagEnvironment), cfg.Env("SERVICE_ENV"), cfg.Const(c.Environment), cfg.Const("development")),
	)
}

var (
	_ cfg.Resolver = (*Config)(nil)
)

func main() {
	flag.Parse()
	config := new(Config)
	if err := cfg.Read(config); !cfg.IsIgnored(err) {
		logger.FatalExit(err)
	}
	fmt.Println("target:", config.Target)
	fmt.Println("env:", config.Environment)
}
