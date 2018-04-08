package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/blend/go-sdk/yaml"
)

var defaultRulesPath = flag.String("f", "./PROFANITY", "the default rules to include for any sub-package")

func main() {
	// walk the filesystem
	// for each file named by the gob filter
	// run the rules on it

	flag.Parse()

	var defaultRules []Rule
	var err error
	if defaultRulesPath != nil && len(*defaultRulesPath) > 0 {
		fmt.Fprintf(os.Stdout, "using default profanity rules file: %s\n", *defaultRulesPath)
		defaultRules, err = deserializeRules(*defaultRulesPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%+v\n", err)
			os.Exit(1)
		}
	}

	packageRules := map[string][]Rule{}

	var getRules = func(path string) ([]Rule, error) {
		if rules, hasRules := packageRules[path]; hasRules {
			return append(defaultRules, rules...), nil
		}
		rules, err := discoverRules(path)
		if err != nil {
			return nil, err
		}
		packageRules[path] = rules
		return append(defaultRules, rules...), nil
	}

	walkErr := filepath.Walk("./", func(file string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() && strings.HasSuffix(file, ".git") {
			return filepath.SkipDir
		}
		if info.IsDir() && strings.HasSuffix(file, "_bin") {
			return filepath.SkipDir
		}

		if !strings.HasSuffix(file, ".go") {
			return nil
		}

		contents, err := ioutil.ReadFile(file)
		if err != nil {
			return err
		}

		rules, err := getRules(filepath.Dir(file))
		if err != nil {
			return err
		}
		for _, rule := range rules {
			if err := rule.Apply(contents); err != nil {
				return fmt.Errorf("%s failed: %+v", file, err)
			}
		}

		return nil
	})

	if walkErr != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", walkErr)
		os.Exit(1)
		return
	}
	fmt.Fprintf(os.Stdout, "profanity ok!\n")
	os.Exit(0)
}

func discoverRules(path string) ([]Rule, error) {
	profanityPath := filepath.Join(path, "PROFANITY")
	if _, err := os.Stat(profanityPath); err != nil {
		return nil, nil
	}
	return deserializeRules(profanityPath)
}

func deserializeRules(path string) (rules []Rule, err error) {
	var contents []byte
	contents, err = ioutil.ReadFile(path)
	if err != nil {
		return
	}

	err = yaml.Unmarshal(contents, &rules)
	return
}

// Contains creates a simple contains rule.
func Contains(value string) RuleFunc {
	return func(contents []byte) error {
		if strings.Contains(string(contents), value) {
			return fmt.Errorf("contains: \"%s\"", value)
		}
		return nil
	}
}

// Regex creates a new regex filter rule.
func Regex(expr string) RuleFunc {
	regex := regexp.MustCompile(expr)
	return func(contents []byte) error {
		if regex.Match(contents) {
			return fmt.Errorf("regexp match: \"%s\"", expr)
		}
		return nil
	}
}

// Rule is a serialized rule.
type Rule struct {
	Contains string `yaml:"contains,omitempty"`
	Regex    string `yaml:"regex,omitempty"`
}

// Apply applies the rule.
func (r Rule) Apply(contents []byte) error {
	if len(r.Contains) > 0 {
		return Contains(r.Contains)(contents)
	}
	if len(r.Regex) > 0 {
		return Regex(r.Regex)(contents)
	}
	return fmt.Errorf("no rule set")
}

// RuleFunc is a function that evaluates a corpus.
type RuleFunc func([]byte) error
