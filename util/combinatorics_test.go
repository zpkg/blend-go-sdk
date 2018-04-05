package util

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestCombinationsOfInt(t *testing.T) {
	assert := assert.New(t)

	combinations := Combinatorics.CombinationsOfInt(1, 2, 3, 4)
	assert.Len(15, combinations)
	assert.Len(4, combinations[0])
}

func TestCombinationsOfFloat(t *testing.T) {
	assert := assert.New(t)

	combinations := Combinatorics.CombinationsOfFloat(1.0, 2.0, 3.0, 4.0)
	assert.Len(15, combinations, 15)
	assert.Len(4, combinations[0])
}

func TestCombinationsOfString(t *testing.T) {
	assert := assert.New(t)

	combinations := Combinatorics.CombinationsOfString("a", "b", "c", "d")
	assert.Len(15, combinations)
	assert.Len(4, combinations[0])
}

func TestPermutationsOfInt(t *testing.T) {
	assert := assert.New(t)

	permutations := Combinatorics.PermutationsOfInt(123, 216, 4, 11)
	assert.Len(24, permutations)
	assert.Len(4, permutations[0])
}

func TestPermutationsOfFloat(t *testing.T) {
	assert := assert.New(t)

	permutations := Combinatorics.PermutationsOfFloat(3.14, 2.57, 1.0, 6.28)
	assert.Len(24, permutations)
	assert.Len(4, permutations[0])
}

func TestPermutationsOfString(t *testing.T) {
	assert := assert.New(t)

	permutations := Combinatorics.PermutationsOfString("a", "b", "c", "d")
	assert.Len(24, permutations)
	assert.Len(4, permutations[0])
}

func TestPermuteDistributions(t *testing.T) {
	assert := assert.New(t)

	permutations := Combinatorics.PermuteDistributions(4, 2)
	assert.Len(5, permutations, 5)
	assert.Len(2, permutations[0])
}

func TestPairsOfInt(t *testing.T) {
	assert := assert.New(t)

	pairs := Combinatorics.PairsOfInt(1, 2, 3, 4, 5)
	assert.Len(10, pairs)
	assert.Equal(4, pairs[9][0])
	assert.Equal(5, pairs[9][1])
}

func TestPairsOfFloat64(t *testing.T) {
	assert := assert.New(t)

	pairs := Combinatorics.PairsOfFloat64(1, 2, 3, 4, 5)
	assert.Len(10, pairs)
	assert.Equal(4, pairs[9][0])
	assert.Equal(5, pairs[9][1])
}

type any = interface{}

func TestAnagrams(t *testing.T) {
	assert := assert.New(t)

	words := Combinatorics.Anagrams("abcde")
	assert.Len(120, words)
	assert.Any(words, func(v any) bool {
		return v.(string) == "abcde"
	})
	assert.Any(words, func(v any) bool {
		return v.(string) == "ecdab"
	})
}
