package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/blend/go-sdk/yaml"
)

// linker metadata block
// this block must be present
// it is used by goreleaser
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

const (
	// DefaultProfanityFile is the default file to use for profanity rules
	DefaultProfanityFile = "PROFANITY"

	// Star is a special character
	Star = "*"
)

var rulesFile = flag.String("rules", DefaultProfanityFile, "the default rules to include for any sub-package.")
var include = flag.String("include", "", "the include file filter in glob form, can be a csv.")
var exclude = flag.String("exclude", "", "the exclude file filter in glob form, can be a csv.")
var verbose = flag.Bool("v", false, "verbose output")

func main() {
	flag.Parse()

	if rulesFile != nil && len(*rulesFile) > 0 {
		if *verbose {
			fmt.Fprintf(os.Stdout, "using rules file: %s\n", *rulesFile)
		}
	}

	if *verbose {
		if len(*include) > 0 {
			fmt.Fprintf(os.Stdout, "using include filter: %s\n", *include)
		}
		if len(*exclude) > 0 {
			fmt.Fprintf(os.Stdout, "using exclude filter: %s\n", *exclude)
		}
	}

	realizedRules := map[string][]Rule{}
	packageRules := map[string][]Rule{}

	var fileBase string
	walkErr := filepath.Walk(".", func(file string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() && strings.HasSuffix(file, ".git") { // don't ever process git directories
			return filepath.SkipDir
		}
		if info.IsDir() {
			return nil
		}

		fileBase = filepath.Base(file)
		if *verbose {
			fmt.Fprintf(os.Stdout, "%s", ColorLightWhite.Apply(file))
		}

		if len(*include) > 0 {
			if matches := globAnyMatch(*include, file); !matches {
				if *verbose {
					fmt.Fprintf(os.Stdout, ".. skipping\n")
				}
				return nil
			}
		}

		if len(*exclude) > 0 {
			if matches := globAnyMatch(*exclude, file); matches {
				if *verbose {
					fmt.Fprintf(os.Stdout, ".. skipping\n")
				}
				return nil
			}
		}

		if matches, err := filepath.Match(DefaultProfanityFile, fileBase); err != nil {
			return err
		} else if matches {
			if *verbose {
				fmt.Fprintf(os.Stdout, ".. skipping\n")
			}
			return nil
		}

		rules, err := getRules(realizedRules, packageRules, filepath.Dir(file))
		if err != nil {
			return err
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

		if *verbose {
			fmt.Fprintf(os.Stdout, " ... %s\n", ColorGreen.Apply("ok!"))
		}

		return nil
	})

	if walkErr != nil {
		fmt.Fprintf(os.Stderr, "%+v\n\n", walkErr)
		os.Exit(1)
		return
	}
	os.Exit(0)
}

func getRules(realizedRules map[string][]Rule, packageRules map[string][]Rule, path string) ([]Rule, error) {
	if rules, hasRules := realizedRules[path]; hasRules {
		return rules, nil
	}

	rules, err := discoverRules(packageRules, path)
	if err != nil {
		return nil, err
	}
	realizedRules[path] = rules
	return rules, nil
}

func discoverRules(packageRules map[string][]Rule, path string) ([]Rule, error) {
	rules, err := localRules(packageRules, path)
	if err != nil {
		return nil, err
	}

	for key, inheritedRules := range packageRules {
		if strings.HasPrefix(path, key) && key != path {
			rules = append(inheritedRules, rules...)
		}
	}

	// always include rules from "." if they were set
	if rootRules, hasRootRules := packageRules["."]; hasRootRules && path != "." {
		rules = append(rootRules, rules...)
	}

	return rules, nil
}

func localRules(packageRules map[string][]Rule, path string) ([]Rule, error) {
	profanityPath := filepath.Join(path, *rulesFile)
	if _, err := os.Stat(profanityPath); err != nil {
		return nil, nil
	}

	rules, err := deserializeRules(profanityPath)
	if err != nil {
		return nil, err
	}
	packageRules[path] = rules
	return rules, nil
}

func deserializeRules(path string) (rules []Rule, err error) {
	var contents []byte
	contents, err = ioutil.ReadFile(path)
	if err != nil {
		return
	}
	var fileRules []Rule
	err = yaml.Unmarshal(contents, &fileRules)
	if err != nil {
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

// Contains creates a simple contains rule.
// It will also return the offending line number.
func Contains(value string) RuleFunc {
	return func(contents []byte) error {
		scanner := bufio.NewScanner(bytes.NewBuffer(contents))
		var line int
		for scanner.Scan() {
			line++
			if strings.Contains(scanner.Text(), value) {
				return fmt.Errorf("contains: \"%s\" (line: %d)", value, line)
			}
		}
		return nil
	}
}

// NotContains creates a simple contains rule.
func NotContains(value string) RuleFunc {
	return func(contents []byte) error {
		if !strings.Contains(string(contents), value) {
			return fmt.Errorf("not contains: \"%s\"", value)
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
	// ID is a unique identifier for the rule.
	ID string `yaml:"id"`
	// File is the rules file path the rule came from.
	File string `yaml:"-"`
	// Message is a descriptive message for the rule.
	Message string `yaml:"message,omitempty"`

	// Include sets a glob filter for file inclusion by filename.
	Include string `yaml:"include,omitempty"`
	// Exclude sets a glob filter for file exclusion by filename.
	Exclude string `yaml:"exclude,omitempty"`

	//
	// the below are matching rules.
	// if these match, the rule is valid
	//

	// Contains implies we should fail if a file contains a given string.
	Contains string `yaml:"contains,omitempty"`
	// Contains implies we should fail if a file doesn't contains a given string.
	NotContains string `yaml:"notContains,omitempty"`
	// Matches implies we should fail if a file's content matches a given regex.
	Matches string `yaml:"matches,omitempty"`
}

// ShouldInclude returns if we should include a file for a given rule.
// If the `.Include` field is unset, this will alway return true.
func (r Rule) ShouldInclude(file string) bool {
	if len(r.Include) == 0 {
		return true
	}
	return globAnyMatch(r.Include, file)
}

// ShouldExclude returns if we should include a file for a given rule.
// If the `.Include` field is unset, this will alway return true.
func (r Rule) ShouldExclude(file string) bool {
	if len(r.Exclude) == 0 {
		return false
	}
	return globAnyMatch(r.Exclude, file)
}

// Apply applies the rule.
func (r Rule) Apply(contents []byte) error {
	if len(r.Contains) > 0 {
		return Contains(r.Contains)(contents)
	}
	if len(r.NotContains) > 0 {
		return NotContains(r.NotContains)(contents)
	}
	if len(r.Matches) > 0 {
		return Regex(r.Matches)(contents)
	}
	return fmt.Errorf("no rule set")
}

// Failure returns a failure error message for a given file and error.
func (r Rule) Failure(file string, err error) error {
	var tokens []string
	if len(r.ID) > 0 {
		tokens = append(tokens, fmt.Sprintf("%s: %s", ColorLightWhite.Apply("rule"), r.ID))
	}

	tokens = append(tokens, fmt.Sprintf("%s %s: %+v", ColorLightWhite.Apply(file), ColorRed.Apply("failed"), err))

	if len(r.Message) > 0 {
		tokens = append(tokens, fmt.Sprintf("%s: %s", ColorLightWhite.Apply("message"), r.Message))
	}
	if len(r.File) > 0 {
		tokens = append(tokens, fmt.Sprintf("%s: %s", ColorLightWhite.Apply("rules file"), r.File))
	}
	if len(r.Include) > 0 {
		tokens = append(tokens, fmt.Sprintf("%s: %s", ColorLightWhite.Apply("include"), r.Include))
	}
	if len(r.Exclude) > 0 {
		tokens = append(tokens, fmt.Sprintf("%s: %s", ColorLightWhite.Apply("exclude"), r.Exclude))
	}

	return fmt.Errorf(strings.Join(tokens, "\n"))
}

// RuleFunc is a function that evaluates a corpus.
type RuleFunc func([]byte) error

// AnsiColor represents an ansi color code fragment.
type AnsiColor string

// escaped escapes the color for use in the terminal.
func (acc AnsiColor) escaped() string {
	return "\033[" + string(acc)
}

// Apply returns a string with the color code applied.
func (acc AnsiColor) Apply(text string) string {
	return acc.escaped() + text + ColorReset.escaped()
}

const (
	// ColorBlack is the posix escape code fragment for black.
	ColorBlack AnsiColor = "30m"

	// ColorRed is the posix escape code fragment for red.
	ColorRed AnsiColor = "31m"

	// ColorGreen is the posix escape code fragment for green.
	ColorGreen AnsiColor = "32m"

	// ColorYellow is the posix escape code fragment for yellow.
	ColorYellow AnsiColor = "33m"

	// ColorBlue is the posix escape code fragment for blue.
	ColorBlue AnsiColor = "34m"

	// ColorPurple is the posix escape code fragement for magenta (purple)
	ColorPurple AnsiColor = "35m"

	// ColorCyan is the posix escape code fragement for cyan.
	ColorCyan AnsiColor = "36m"

	// ColorWhite is the posix escape code fragment for white.
	ColorWhite AnsiColor = "37m"

	// ColorLightBlack is the posix escape code fragment for black.
	ColorLightBlack AnsiColor = "90m"

	// ColorLightRed is the posix escape code fragment for red.
	ColorLightRed AnsiColor = "91m"

	// ColorLightGreen is the posix escape code fragment for green.
	ColorLightGreen AnsiColor = "92m"

	// ColorLightYellow is the posix escape code fragment for yellow.
	ColorLightYellow AnsiColor = "93m"

	// ColorLightBlue is the posix escape code fragment for blue.
	ColorLightBlue AnsiColor = "94m"

	// ColorLightPurple is the posix escape code fragement for magenta (purple)
	ColorLightPurple AnsiColor = "95m"

	// ColorLightCyan is the posix escape code fragement for cyan.
	ColorLightCyan AnsiColor = "96m"

	// ColorLightWhite is the posix escape code fragment for white.
	ColorLightWhite AnsiColor = "97m"

	// ColorGray is an alias to ColorLightWhite to preserve backwards compatibility.
	ColorGray AnsiColor = ColorLightBlack

	// ColorReset is the posix escape code fragment to reset all formatting.
	ColorReset AnsiColor = "0m"
)

// globIncludeMatch tests if a file matches a (potentially) csv of glob filters.
func globAnyMatch(filter, file string) bool {
	parts := strings.Split(filter, ",")
	for _, part := range parts {
		if matches := glob(strings.TrimSpace(part), file); matches {
			return true
		}
	}
	return false
}

func glob(pattern, subj string) bool {
	// Empty pattern can only match empty subject
	if pattern == "" {
		return subj == pattern
	}

	// If the pattern _is_ a glob, it matches everything
	if pattern == Star {
		return true
	}

	parts := strings.Split(pattern, Star)

	if len(parts) == 1 {
		// No globs in pattern, so test for equality
		return subj == pattern
	}

	leadingGlob := strings.HasPrefix(pattern, Star)
	trailingGlob := strings.HasSuffix(pattern, Star)
	end := len(parts) - 1

	// Go over the leading parts and ensure they match.
	for i := 0; i < end; i++ {
		idx := strings.Index(subj, parts[i])

		switch i {
		case 0:
			// Check the first section. Requires special handling.
			if !leadingGlob && idx != 0 {
				return false
			}
		default:
			// Check that the middle parts match.
			if idx < 0 {
				return false
			}
		}

		// Trim evaluated text from subj as we loop over the pattern.
		subj = subj[idx+len(parts[i]):]
	}

	// Reached the last section. Requires special handling.
	return trailingGlob || strings.HasSuffix(subj, parts[end])
}
