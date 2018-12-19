package stringutil

import (
	"net/url"
	"unicode"
)

// Slugify replaces whitespace with '-' and url escapes.
func Slugify(v string) string {
	runes := []rune(v)
	var c rune
	for index := range runes {
		c = runes[index]
		if unicode.IsSpace(c) {
			runes[index] = '-'
		}
	}
	return url.PathEscape(string(runes))
}
