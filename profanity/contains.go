package profanity

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"
)

// Contains creates a simple contains rule.
// It fails if a corpus contains a given value.
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
