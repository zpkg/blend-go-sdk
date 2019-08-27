package env

import (
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

// PairDelimiter is a type of delimiter that separates different env var key-value pairs
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
			serializedPair := QuoteDelimiter + escapeString(k, separator) +
				QuoteDelimiter + ValueDelimiter + QuoteDelimiter +
				escapeString(v, separator) + QuoteDelimiter + separator
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

	// indicates whether the value delimiter has been encountered for the current pair
	valueFlag := false

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
				if _, exists := ret[key]; exists {
					return ret, ex.New("Duplicate keys are not allowed")
				}

				if len(key) == 0 {
					return ret, ex.New("Empty keys are not allowed")
				}

				// This means that we have a term with no '=', which is illegal
				if !valueFlag {
					return ret, ex.New("Expected '='")
				}

				value = buffer
				ret[key] = value

				// clear out the buffers and start over
				buffer = ""
				key = ""
				value = ""
				valueFlag = false
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
			// State 1: escape literal -- we want to take whatever the next
			// token is no matter what, goes back to the root mode
			buffer += char
			state = 0
		case 2:
			// State 2: process the '=' character. We need to reset the buffer,
			// store the key, and start storing characters in the buffer that
			// will go to the value
			if len(buffer) == 0 {
				return ret, ex.New("Empty keys are not allowed")
			}
			key = buffer
			buffer = ""
			valueFlag = true

			if char == QuoteDelimiter {
				state = 3
			} else {
				if !unicode.IsSpace(c) {
					buffer += char
				}
				state = 0
			}
		case 3:
			// State 3: quote mode -- accept all text except for the end quote
			// (excluding anything that is escaped)
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

	// State 0 is the only valid ending state. If this is not the case, then
	// show the user a parsing error. In the event the input wasn't terminated,
	// we can mitigate by taking the last key-val pair from the buffers.
	switch state {
	case 0:
		// This handles the case where the key-value pair doesn't have a
		// separator (which is valid grammar). We could go about the option of
		// inserting an extra separator, but that is difficult to do as a
		// preprocessing step because you could have a scenario where there are
		// trailing spaces, or even an escaped ending delimiter.
		if len(buffer) > 0 || len(key) > 0 {
			if !valueFlag {
				return ret, ex.New("Expected '='")
			}
			ret[key] = buffer
		}
	case 1:
		return ret, ex.New("Ended input on an escape delimiter ('\\')")
	case 2:
		return ret, ex.New("Failed to assign a value to some key")
	case 3:
		return ret, ex.New("Unclosed quote")
	case 4:
		return ret, ex.New("Ended input on an escape delimiter ('\\')")
	}
	return ret, nil
}

// isToken returns whether a string is a special token that would need to be
// escaped
func isToken(s string, delimiter PairDelimiter) bool {
	switch s {
	case delimiter,
		ValueDelimiter,
		QuoteDelimiter,
		EscapeDelimiter:
		return true
	}
	return false
}

// escapeString takes a string and escapes any special characters so that the
// string can be serialized properly. The user must supply the delimiter used
// separate key-value pairs.
func escapeString(s string, delimiter PairDelimiter) string {
	var escaped string

	for _, r := range s {
		char := string(r)

		if isToken(char, delimiter) {
			escaped += EscapeDelimiter
		}
		escaped += char
	}
	return escaped
}
