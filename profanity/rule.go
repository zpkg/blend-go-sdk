package profanity

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/blend/go-sdk/ansi"
	"github.com/blend/go-sdk/exception"
	"github.com/blend/go-sdk/stringutil"
	"github.com/blend/go-sdk/yaml"
)

// Rules are a list of rules.
type Rules []Rule

// String
func (r Rules) String() string {
	var output string
	for _, rule := range r {
		output = output + rule.String() + "\n"
	}
	return output
}

// RulesFromPath reads rules from a path
func RulesFromPath(path string) (rules []Rule, err error) {
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

// Rule is a serialized rule.
type Rule struct {
	// ID is a unique identifier for the rule.
	ID string `yaml:"id"`
	// File is the rules file path the rule came from.
	File string `yaml:"-"`
	// Message is a descriptive message for the rule.
	Message string `yaml:"message,omitempty"`

	// IncludeAny sets a glob filter for file inclusion by filename.
	IncludeAny []string `yaml:"includeAny,omitempty"`
	// ExcludeAny sets a glob filter for file exclusion by filename.
	ExcludeAny []string `yaml:"excludeAny,omitempty"`

	//
	// the below are matching rules.
	// if these match, the rule will fail the profanity check
	//

	// ContainsAny implies we should fail if a file contains a given string.
	ContainsAny []string `yaml:"containsAny,omitempty"`
	// NotContainsAll implies we should fail if a file doesn't contains a given string.
	NotContainsAll []string `yaml:"notContainsAll,omitempty"`
	// Matches implies we should fail if a file's content matches a given regex.
	MatchesAny []string `yaml:"matchesAny,omitempty"`
}

// ShouldInclude returns if we should include a file for a given rule.
// If the `.Include` field is unset, this will alway return true.
func (r Rule) ShouldInclude(file string) bool {
	if len(r.IncludeAny) == 0 {
		return true
	}
	return GlobAnyMatch(r.IncludeAny, file)
}

// ShouldExclude returns if we should include a file for a given rule.
// If the `.Include` field is unset, this will alway return true.
func (r Rule) ShouldExclude(file string) bool {
	if len(r.ExcludeAny) == 0 {
		return false
	}
	return GlobAnyMatch(r.ExcludeAny, file)
}

// Apply applies the rule.
func (r Rule) Apply(contents []byte) error {
	if len(r.ContainsAny) > 0 {
		return ContainsAny(r.ContainsAny...)(contents)
	}
	if len(r.NotContainsAll) > 0 {
		return NotContainsAll(r.NotContainsAll...)(contents)
	}
	if len(r.MatchesAny) > 0 {
		return MatchesAny(r.MatchesAny...)(contents)
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
	if len(r.IncludeAny) > 0 {
		tokens = append(tokens, fmt.Sprintf("%s: %s", ansi.LightWhite("includes"), stringutil.CSV(r.IncludeAny)))
	}
	if len(r.ExcludeAny) > 0 {
		tokens = append(tokens, fmt.Sprintf("%s: %s", ansi.LightWhite("excludes"), stringutil.CSV(r.ExcludeAny)))
	}

	return fmt.Errorf(strings.Join(tokens, "\n"))
}

// String returns a string representation of the rule.
func (r Rule) String() string {
	var tokens []string

	if len(r.ID) > 0 {
		tokens = append(tokens, fmt.Sprintf("[%s]", r.ID))
	}
	if len(r.Message) > 0 {
		tokens = append(tokens, "`"+r.Message+"`")
	}
	if len(r.IncludeAny) > 0 {
		tokens = append(tokens, fmt.Sprintf("[include any: %s]", strings.Join(r.IncludeAny, ", ")))
	}
	if len(r.ExcludeAny) > 0 {
		tokens = append(tokens, fmt.Sprintf("[exclude any: %s]", strings.Join(r.ExcludeAny, ",")))
	}
	if len(r.ContainsAny) > 0 {
		tokens = append(tokens, fmt.Sprintf("[contains any: %s]", strings.Join(r.ContainsAny, ",")))
	}
	if len(r.NotContainsAll) > 0 {
		tokens = append(tokens, fmt.Sprintf("[not contains all: %s]", strings.Join(r.NotContainsAll, ",")))
	}
	if len(r.MatchesAny) > 0 {
		tokens = append(tokens, fmt.Sprintf("[matches any: %s]", strings.Join(r.MatchesAny, ",")))
	}
	return strings.Join(tokens, " ")
}
