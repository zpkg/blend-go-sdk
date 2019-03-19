package profanity

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"

	"github.com/blend/go-sdk/ansi"
	"github.com/blend/go-sdk/yaml"
)

// RulesFromPath reads rules from a path
func RulesFromPath(path string) (rules []Rule, err error) {
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
	// if these match, the rule will fail the profanity check
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
	return GlobAnyMatch(r.Include, file)
}

// ShouldExclude returns if we should include a file for a given rule.
// If the `.Include` field is unset, this will alway return true.
func (r Rule) ShouldExclude(file string) bool {
	if len(r.Exclude) == 0 {
		return false
	}
	return GlobAnyMatch(r.Exclude, file)
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
		tokens = append(tokens, fmt.Sprintf("%s: %s", ansi.LightWhite("rule"), r.ID))
	}

	tokens = append(tokens, fmt.Sprintf("%s %s: %+v", ansi.LightWhite(file), ansi.Red("failed"), err))

	if len(r.Message) > 0 {
		tokens = append(tokens, fmt.Sprintf("%s: %s", ansi.LightWhite("message"), r.Message))
	}
	if len(r.File) > 0 {
		tokens = append(tokens, fmt.Sprintf("%s: %s", ansi.LightWhite("rules file"), r.File))
	}
	if len(r.Include) > 0 {
		tokens = append(tokens, fmt.Sprintf("%s: %s", ansi.LightWhite("include"), r.Include))
	}
	if len(r.Exclude) > 0 {
		tokens = append(tokens, fmt.Sprintf("%s: %s", ansi.LightWhite("exclude"), r.Exclude))
	}

	return fmt.Errorf(strings.Join(tokens, "\n"))
}
