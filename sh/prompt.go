package sh

import (
	"fmt"
	"github.com/blend/go-sdk/exception"
)

const ErrUnexpectedNewline exception.Class = "unexpected newline"

// MustPrompt gives a prompt and reads input until newlines.
// It panics on error.
func MustPrompt(prompt string) string {
	output, err := Prompt(prompt)
	if err != nil {
		if exception.As(err).Class() == ErrUnexpectedNewline {
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
	if exception.As(err).Class() == ErrUnexpectedNewline {
		return "", nil
	}
	return output, err
}
