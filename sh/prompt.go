package sh

import "fmt"

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
	fmt.Print(prompt)
	var output string
	_, err := fmt.Scanln(&output)
	return output, err
}
