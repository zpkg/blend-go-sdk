package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/blend/go-sdk/configutil"
	"github.com/blend/go-sdk/profanity"
	"github.com/blend/go-sdk/ref"
)

// linker metadata block
// this block must be present
// it is used by goreleaser
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

var (
	flagRulesFile = flag.String("rules", "", "the default rules to include for any sub-package.")
	flagInclude   = flag.String("include", "", "the include file filter in glob form, can be a csv.")
	flagExclude   = flag.String("exclude", "", "the exclude file filter in glob form, can be a csv.")
	flagVerbose   = flag.Bool("v", false, "verbose output")
)

type config struct {
	profanity.Config `yaml:",inline"`
}

// Resolve resolves the config.
func (c *config) Resolve() error {
	return configutil.AnyError(
		configutil.SetBool(&c.Verbose, configutil.Bool(flagVerbose), configutil.Bool(c.Verbose), configutil.Bool(ref.Bool(false))),
		configutil.SetString(&c.RulesFile, configutil.String(*flagRulesFile), configutil.String(c.RulesFile)),
		configutil.SetString(&c.Include, configutil.String(*flagInclude), configutil.String(c.Include)),
		configutil.SetString(&c.Exclude, configutil.String(*flagExclude), configutil.String(c.Exclude)),
	)
}

func main() {
	flag.Parse()

	var cfg config
	if err := cfg.Resolve(); err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n\n", err)
		os.Exit(1)
	}

	engine := profanity.New(profanity.OptConfig(&cfg.Config))
	if err := engine.Process(); err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n\n", err)
		os.Exit(1)
		return
	}
	os.Exit(0)
}
