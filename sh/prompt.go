package sh

import (
	"fmt"
	"io"
	"os"

	"github.com/blend/go-sdk/exception"
)

// ErrUnexpectedNewLine is returned from scan.go when you just hit enter with nothing in the prompt
const ErrUnexpectedNewLine exception.Class = "unexpected newline"

// MustPrompt gives a prompt and reads input until newlines.
// It panics on error.
func MustPrompt(prompt string) string {
	output, err := Prompt(prompt)
	if err != nil {
		panic(err)
	}
	return output
}

// Prompt gives a prompt and reads input until newlines.
func Prompt(prompt string) (string, error) {
	return PromptFrom(os.Stdout, os.Stdin, prompt)
}

// Promptf gives a prompt of a given format and args and reads input until newlines.
func Promptf(format string, args ...interface{}) (string, error) {
	return PromptFrom(os.Stdout, os.Stdin, fmt.Sprintf(format, args...))
}

// PromptFrom gives a prompt and reads input until newlines from a given set of streams.
func PromptFrom(stdout io.Writer, stdin io.Reader, prompt string) (string, error) {
	fmt.Fprint(stdout, prompt)
	var output string
	_, err := fmt.Fscanln(stdin, &output)
	if exception.Is(ErrUnexpectedNewLine, err) {
		return "", nil
	}
	return output, err
}
