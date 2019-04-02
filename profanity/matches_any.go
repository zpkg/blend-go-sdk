package profanity

import (
	"fmt"
	"regexp"
)

// MatchesAny creates a new regex filter rule.
// It failes if any of the expressions match.
func MatchesAny(exprs ...string) RuleFunc {
	return func(contents []byte) error {
		for _, expr := range exprs {
			regex := regexp.MustCompile(expr)
			if regex.Match(contents) {
				return fmt.Errorf("regexp match: \"%s\"", expr)
			}
		}
		return nil
	}
}
