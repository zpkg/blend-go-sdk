package sh

import (
	"fmt"
)

// ErrUnexpectedNewline is returned from scan.go when you just hit enter with nothing in the prompt
const ErrUnexpectedNewline string = "unexpected newline"

// MustPrompt gives a prompt and reads input until newlines.
// It panics on error.
func MustPrompt(prompt string) string {
	output, err := Prompt(prompt)
	if err != nil {
		if err.Error() == ErrUnexpectedNewline {
			return ""
		}
		panic(err)
	}
	return output
}

// Prompt gives a prompt and reads input until newlines.
func Prompt(prompt string) (string, error) {
	fmt.Print(prompt)
	var output string
	_, err := fmt.Scanln(&output)
	if err.Error() == ErrUnexpectedNewline {
		return "", nil
	}
	return output, err
}
