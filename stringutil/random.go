package stringutil

import (
	"math/rand"
	"time"
)

var (
	provider = rand.New(rand.NewSource(time.Now().UnixNano()))
)

// Random returns a random selection of runes from the set.
//
// NOTE: This should not be used in security related settings because the random source
// is not guaranteed to be secure. Instead use `crypto.CreateKey(...)` if you need to generate
// a secure random value.
func Random(runeset []rune, length int) string {
	return Runeset(runeset).Random(length)
}
