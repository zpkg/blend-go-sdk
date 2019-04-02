package profanity

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"
)

// ContainsAny creates a simple contains rule.
// It acts as an OR; it fails if a corpus contains any given value.
func ContainsAny(values ...string) RuleFunc {
	return func(filename string, contents []byte) error {
		scanner := bufio.NewScanner(bytes.NewBuffer(contents))
		var line int
		for scanner.Scan() {
			line++
			for _, value := range values {
				if strings.Contains(scanner.Text(), value) {
					return fmt.Errorf("contains: \"%s\" (line: %d)", value, line)
				}
			}
		}
		return nil
	}
}
