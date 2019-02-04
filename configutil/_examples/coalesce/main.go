package main

import (
	"flag"
	"fmt"

	"github.com/blend/go-sdk/configutil"
	"github.com/blend/go-sdk/logger"
)

// Config is a sample config.
type Config struct {
	Target string `yaml:"target"`
}

// Resolve parses the config and sets values from predefined sources.
func (c *Config) Resolve() error {
	if err := configutil.CoalesceSourcesVar(&c.Target,
		configutil.FlagSource("target", "The intended target"),
		configutil.EnvSource("TARGET"),
		configutil.ConstantSource(c.Target),
		configutil.ConstantSource("https://google.com/robots.txt"),
	); err != nil {
		return err
	}
	return nil
}

func main() {
	config := new(Config)
	if err := configutil.Read(config); !configutil.IsIgnored(err) {
		logger.FatalExit(err)
	}
	if err := config.Resolve(flag.CommandLine); err != nil {
		logger.FatalExit(err)
	}
	fmt.Println("target:", config.Target)
}
