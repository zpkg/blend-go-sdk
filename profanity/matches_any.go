package profanity

import (
	"bufio"
	"bytes"
	"fmt"
	"regexp"
)

// MatchesAny creates a new regex filter rule.
// It failes if any of the expressions match.
func MatchesAny(exprs ...string) RuleFunc {
	return func(filename string, contents []byte) RuleResult {
		scanner := bufio.NewScanner(bytes.NewBuffer(contents))
		var line int
		for scanner.Scan() {
			line++
			for _, expr := range exprs {
				regex := regexp.MustCompile(expr)
				if regex.Match([]byte(scanner.Text())) {
					return RuleResult{
						File:    filename,
						Line:    line,
						Message: fmt.Sprintf("regexp match: \"%s\"", expr),
					}
				}
			}
		}
		return RuleResult{OK: true}
	}
}
