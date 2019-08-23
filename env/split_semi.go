package env

import (
	"fmt"
	"unicode"

	"github.com/blend/go-sdk/ex"
)

const (
	// ValueDelimiter ("=") is the delimiter between a key and a value for an
	// environment variable.
	ValueDelimiter = "="

	// QuoteDelimiter (`"`) is a delimiter indicating a string literal. This
	// gives the user the option to have spaces, for example, in their
	// environment variable values.
	QuoteDelimiter = "\""

	// EscapeDelimiter ("\") is used to escape the next character so it is
	// accepted as a part of the input value.
	EscapeDelimiter = "\\"

	// SpaceDelimiter (" ") is a delimiter that simply represents a space. It
	// is generally ignored, unless quoted.
	SpaceDelimiter = " "
)

// PairDelimiter is the type of delimiter that separates different env var key-value pairs
type PairDelimiter = string

const (
	// SemicolonDelimiter (";") is a delimiter between key-value pairs
	SemicolonDelimiter PairDelimiter = ";"

	// CommaDelimiter (",") is a delimiter betewen key-value pairs
	CommaDelimiter PairDelimiter = ","
)

// delimitedString converts environment variables to a particular string
// representation, allowing the user to specify which delimiter to use between
// different environment variable pairs.
func (ev Vars) DelimitedString(separator PairDelimiter) string {
	res := ""

	// For each key, value pair, convert it into a "key=value;" pair and
	// continue appending to the output string for each pair
	for k, v := range ev {
		if k != "" {
			serializedPair := fmt.Sprintf("%s=%s%s", k, v, separator)
			res += serializedPair
		}
	}
	return res
}

// Parse uses a state machine to parse an input string into the `Vars` type.
// The user can choose which delimiter to use between key-value pairs.
//
// An example of this format:
//
// ENV_VAR_1=VALUE_1;ENV_VAR_2=VALUE_2;
//
// We define the grammar as such (in BNF notation):
// <expr> ::= (<pair> <sep>)* <pair>
// <sep> ::= ';'
//        |  ','
// <pair> ::= <term> = <term>
// <term> ::= <literal>
//         |  "[<literal>|<space>|<escape_quote>]*"
// <literal> ::= [-A-Za-z_0-9]+
// <space> ::= ' '
// <escape_quote> ::= '\"'
func Parse(s string, separator PairDelimiter) (Vars, error) {
	ret := make(Vars)
	var key string
	var value string
	var buffer string
	state := 0

	for _, c := range s {
		// Having a string is convenient so we can do quality comparisons with
		// the tokens defined in the constants section
		char := string(c)

		switch state {
		// The "root" case, which simply evaluates each character from the
		// initial state. This is the only valid ending state.
		case 0:
			// In the case where we have a key-value pair, we want to add that
			// to the map and clear out our buffers
			if char == separator {
				// check that we don't have a duplicate
				if _, exists := ret[key]; exists {
					return ret, ex.New("Duplicate keys are not allowed")
				}

				// It is illegal to have an empty key
				if len(key) == 0 {
					return ret, ex.New("Empty keys are not allowed")
				}

				// It is illegal to have an empty value
				if len(buffer) == 0 {
					return ret, ex.New("Empty values are not allowed")
				}
				value = buffer
				ret[key] = value

				// clear out the buffers and start over
				buffer = ""
				key = ""
				value = ""
			} else if char == EscapeDelimiter {
				state = 1
			} else if char == ValueDelimiter {
				state = 2
			} else if char == QuoteDelimiter {
				state = 3
			} else if unicode.IsSpace(c) {
				continue
			} else {
				buffer += char
			}
		case 1:
			// Escape literal: we want to take whatever the next token is, goes back to the root mode
			buffer += char
			state = 0
		case 2:
			// It is illegal to have an empty key
			if len(buffer) == 0 {
				return ret, ex.New("Empty keys are not allowed")
			}
			// Previous value was an equals sign, so we want to assign the key and
			// clear the buffer and add the current character
			key = buffer
			// clear the buffer and add the current character (excluding whitespace)
			buffer = ""

			if !unicode.IsSpace(c) {
				buffer += char
			}
			state = 0
		case 3:
			// Quote mode: accept all text except for the end quote (excluding anything that is escaped)
			if char == EscapeDelimiter {
				// ignore the escape and continue
				state = 4
			} else if char == QuoteDelimiter {
				// go back to the default state
				state = 0
			} else {
				buffer += char
			}
		case 4:
			// Escape literal within a quote, goes back to quote mode
			buffer += char
			state = 3
		}
	}

	// This handles the case where the key-value pair doesn't have a separator
	// (which is valid grammar). We could go about the option of inserting an
	// extra separator, but that is difficult to do as a preprocessing step
	// because you could have a scenario where there are trailing spaces, or even
	if state == 0 && len(buffer) > 0 {
		ret[key] = buffer
	}

	// State 0 is the only valid ending state. If this is not the case, then
	// show the user a parsing error. In the event the input wasn't terminated,
	// we can mitigate by taking the last key-val pair from the buffers.
	switch state {
	case 1:
		return ret, ex.New("Ended input on an escape delimiter (`\\`)")
	case 2:
		return ret, ex.New("Failed to assign a value to some key")
	case 3:
		return ret, ex.New("Unclosed quote")
	case 4:
		return ret, ex.New("Ended input on an escape delimiter (`\\`)")
	}
	return ret, nil
}
