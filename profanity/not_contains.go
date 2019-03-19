package profanity

import (
	"fmt"
	"strings"
)

// NotContains creates a simple not contains rule.
// It fails if a corpus does not contain a given value.
func NotContains(value string) RuleFunc {
	return func(contents []byte) error {
		if !strings.Contains(string(contents), value) {
			return fmt.Errorf("not contains: \"%s\"", value)
		}
		return nil
	}
}
