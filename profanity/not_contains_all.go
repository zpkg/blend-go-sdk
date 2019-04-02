package profanity

import (
	"fmt"
	"strings"
)

// NotContainsAll creates a simple not contains rule.
// It fails if a corpus does not contain a given value.
// These values act as an AND, as in, they must all be present.
func NotContainsAll(values ...string) RuleFunc {
	return func(contents []byte) error {
		for _, value := range values {
			if !strings.Contains(string(contents), value) {
				return fmt.Errorf("not contains: \"%s\"", value)
			}
		}
		return nil
	}
}
