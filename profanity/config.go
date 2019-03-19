package profanity

// Config is the profanity rules parsing config.
type Config struct {
	Verbose   *bool  `yaml:"verbose"`
	Root      string `yaml:"root"`
	RulesFile string `yaml:"rulesFile"`
	Include   string `yaml:"include"`
	Exclude   string `yaml:"exclude"`
}

// VerboseOrDefault returns an option or a default.
func (c Config) VerboseOrDefault() bool {
	if c.Verbose != nil {
		return *c.Verbose
	}
	return false
}

// RulesFileOrDefault returns the rules file or a default.
func (c Config) RulesFileOrDefault() string {
	if c.RulesFile != "" {
		return c.RulesFile
	}
	return DefaultRulesFile
}
