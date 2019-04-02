package profanity

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/blend/go-sdk/ansi"
	"github.com/blend/go-sdk/exception"
	"github.com/blend/go-sdk/stringutil"
	"github.com/blend/go-sdk/yaml"
)

// New creates a new profanity engine with a given set of config options.
func New(options ...ConfigOption) *Profanity {
	var cfg Config
	for _, option := range options {
		option(&cfg)
	}
	return &Profanity{
		Config: &cfg,
	}
}

// Profanity parses rules from the filesystem and applies them to a given root path.
// Creating a full rules set.
type Profanity struct {
	Config *Config
	Stdout io.Writer
	Stderr io.Writer
}

// Printf writes to the output stream.
func (p *Profanity) Printf(format string, args ...interface{}) {
	if p.Stdout != nil {
		fmt.Fprintf(p.Stdout, format, args...)
	}
}

// Errorf writes to the error output stream.
func (p *Profanity) Errorf(format string, args ...interface{}) {
	if p.Stdout != nil {
		fmt.Fprintf(p.Stderr, format, args...)
	}
}

// Process processes the profanity rules.
func (p *Profanity) Process() error {
	if p.Config.VerboseOrDefault() {
		p.Printf("using rules file: %s\n", p.Config.RulesFileOrDefault())
	}

	if p.Config.VerboseOrDefault() {
		if len(p.Config.Include) > 0 {
			p.Printf("using include filter: %s\n", p.Config.Include)
		}
		if len(p.Config.Exclude) > 0 {
			p.Printf("using exclude filter: %s\n", p.Config.Exclude)
		}
	}

	ruleCache := map[string][]Rule{}

	var fileBase string
	return filepath.Walk(".", func(file string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() && strings.HasSuffix(file, ".git") { // don't ever process git directories
			if p.Config.VerboseOrDefault() {
				p.Printf("%s ... skipping (is .git dir)\n", ansi.LightWhite(file))
			}
			return filepath.SkipDir
		}
		if info.IsDir() {
			return nil
		}

		fileBase = filepath.Base(file)

		if len(p.Config.Include) > 0 {
			if matches := GlobAnyMatch(p.Config.Include, file); !matches {
				if p.Config.VerboseOrDefault() {
					p.Printf("%s ... skipping (does not match include filter)\n", ansi.LightWhite(file))
				}
				return nil
			}
		}

		if len(p.Config.Exclude) > 0 {
			if matches := GlobAnyMatch(p.Config.Exclude, file); matches {
				if p.Config.VerboseOrDefault() {
					p.Printf("%s ... skipping (matches exclude filter)\n", ansi.LightWhite(file))
				}
				return nil
			}
		}

		if matches, err := filepath.Match(p.Config.RulesFileOrDefault(), fileBase); err != nil {
			return err
		} else if matches {
			if p.Config.VerboseOrDefault() {
				p.Printf("%s ... skipping (is %s file)\n", ansi.LightWhite(file), p.Config.RulesFileOrDefault())
			}
			return nil
		}

		fullPath := filepath.Dir(file)
		rules, err := p.RulesForPathOrCached(ruleCache, fullPath)
		if err != nil {
			return err
		}

		if p.Config.VerboseOrDefault() {
			p.Printf("%s rules:\n%s", ansi.LightWhite(file), stringutil.Indent("\t", Rules(rules).String()))
		}

		contents, err := ioutil.ReadFile(file)
		if err != nil {
			return err
		}

		for _, rule := range rules {
			if matches := rule.ShouldInclude(file); !matches {
				continue
			}

			if matches := rule.ShouldExclude(file); matches {
				continue
			}

			if err := rule.Apply(contents); err != nil {
				return rule.Failure(file, err)
			}
		}

		if p.Config.VerboseOrDefault() {
			p.Printf("%s ... %s\n", ansi.LightWhite(file), ansi.Green("ok!"))
		}

		return nil
	})
}

// RulesForPathOrCached returns rules cached or rules from disk.
// It prevents re-reading the full rules set for each file in a path.
func (p *Profanity) RulesForPathOrCached(packageRules map[string][]Rule, path string) ([]Rule, error) {
	if rules, ok := packageRules[path]; ok {
		return rules, nil
	}

	rules, err := p.RulesForPath(packageRules, path)
	if err != nil {
		return nil, err
	}
	packageRules[path] = rules
	return rules, nil
}

// RulesForPath adds rules in a given path and child paths to an existing rule set.
// `workingSet` are the current working rules keyed on the path they
// came from, including '.' for the root rules.
func (p *Profanity) RulesForPath(workingSet map[string][]Rule, path string) ([]Rule, error) {
	pathRules, err := p.ReadRules(path)
	if err != nil {
		return nil, err
	}

	for key, workingRules := range workingSet {
		if strings.HasPrefix(path, key) && key != path {
			pathRules = append(workingRules, pathRules...)
		}
	}

	// always include rules from "." if they were set
	if rootRules, hasRootRules := workingSet[Root]; hasRootRules && path != Root {
		pathRules = append(rootRules, pathRules...)
	}

	return pathRules, nil
}

// ReadRules reads rules at a given directory path.
// Path is meant to be the slash terminated dir, which will have the configured rule path appended to it.
func (p *Profanity) ReadRules(path string) ([]Rule, error) {
	profanityPath := filepath.Join(path, p.Config.RulesFileOrDefault())
	if _, err := os.Stat(profanityPath); err != nil {
		return nil, nil
	}

	if p.Config.VerboseOrDefault() {
		p.Printf("%s reading rules file %s\n", ansi.LightWhite(path), p.Config.RulesFileOrDefault())
	}
	rules, err := p.RulesFromPath(profanityPath)
	if err != nil {
		return nil, err
	}
	return rules, nil
}

// RulesFromPath reads rules from a path
func (p *Profanity) RulesFromPath(path string) (rules []Rule, err error) {
	var contents []byte
	contents, err = ioutil.ReadFile(path)
	if err != nil {
		err = exception.New(err, exception.OptMessagef("file: %s", path))
		return
	}
	var fileRules []Rule
	err = yaml.Unmarshal(contents, &fileRules)
	if err != nil {
		err = exception.New(err, exception.OptMessagef("file: %s", path))
		return
	}
	rules = make([]Rule, len(fileRules))
	for index, fileRule := range fileRules {
		rule := fileRule
		rule.File = path
		rules[index] = rule
	}
	return
}
