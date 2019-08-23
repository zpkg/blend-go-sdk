package env

import (
	"fmt"
	"strings"
)

const (
	// ValueDelimiter ("=") is the delimiter between a key and a value for an
	// environment variable.
	ValueDelimiter = "="

	// QuoteDelimiter (`"`) is a delimiter indicating a string literal. This
	// gives the user the option to have spaces, for example, in their
	// environment variable values.
	QuoteDelimiter = "\""
)

// PairDelimiter is the type of delimiter that separates different env var key-value pairs
type PairDelimiter = string

const (
	// SemicolonDelimiter (";") is a delimiter between key-value pairs
	SemicolonDelimiter PairDelimiter = ";"

	// CommaDelimiter (",") is a delimiter betewen key-value pairs
	CommaDelimiter PairDelimiter = ","
)

// Parse environment variables from from semicolon delimited strings
// We define semicolon delimited as: "KEY_1=VAL;KEY_2=VAL;..."
// TODO(afnan) check to ensure that an empty input returns
func FromSemi(input string) Vars {
	vars := make(Vars)

	// Tokenize by semicolon to separate key-value pairs
	keyVals := strings.Split(input, ";")

	// Tokenize by '=' to separate keys and values.
	for _, pair := range keyVals {
		// We use `SplitN` so that people can have env vars with a '=' in the
		// value, otherwise there will be issues with any value that contains '='
		tokenizedPair := strings.SplitN(pair, "=", 2)
		key := tokenizedPair[0]

		// set the default value just in case the input is deformed -- if for
		// some reason we have no actual value, this will at least indicate
		// that there was some environment variable set
		value := ""

		if len(tokenizedPair) > 1 {
			value = tokenizedPair[1]
		}
		vars[key] = value
	}
	return vars
}

// delimitedString converts environment variables to a particular string
// representation, allowing the user to specify which delimiter to use between
// different environment variable pairs.
func (ev Vars) delimitedString(separator PairDelimiter) string {
	res := ""

	// For each key, value pair, convert it into a "key=value;" pair and
	// continue appending to the output string for each pair
	for k, v := range ev {
		if k != "" {
			serializedPair := fmt.Sprintf("%s%s%s;", k, separator, v)
			res += serializedPair
		}
	}
	return res
}

// Parse uses a state machine to parse an input string into the `Vars` type.
// The user can choose which delimiter to use between key-value pairs. This is
// an LL(1) recursive descent parser.
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
//
// We also define a DFA, which corresponds to the switch-cases in the parsing
// routine.
func Parse(input string, separator PairDelimiter) (Vars, error) {
	// TODO(afnan)
	var ret Vars

	// If the last character isn't a delimiter, inject it at the end to make
	// parsing easier.
	if input[len(input)-1:] != separator {
		input = input + separator
	}
	return ret, nil
}

// expr is a subroutine for the `Parse` method
func expr(s string, separator PairDelimiter) (Vars, error) {
	var ret Vars
	var key string
	var value string
	var literal string
	//startIdx := 0
	state := 0

	for c := range s {
		char := string(c)

		switch state {
		// the "root" case, which simply evaluates each character from the initial state
		case 0:
			if char == separator {
				ret[key] = value
			} else if char == "\\" {
				state = 1
			}

			literal += char
			// with an escape literal, we want to take WHATEVER the next token is
		case 1:
			literal += char
		}
	}
	for _, char := range s {
		charString := string(char)
		if charString == separator {

		}
	}
	return ret, nil
}

// pair is a subroutine for the `Parse` method. This function will return an
// environment variable key-value pair.
func pair(s string) (string, string, error) {
	// return an environment variable pair
	return "", "", nil
}

// term is a subroutine for the `Parse` method that determines whether the
// input is parsing a valid term, with optional quotes.
func term(s string) (string, error) {
	return "", nil
}

// literal is a subroutine for the `Parse` method. This function determines
// whether a character is a valid literal.
// TODO(afnan) remove since this is a noop
func literal(s string) bool {
	return true
}
