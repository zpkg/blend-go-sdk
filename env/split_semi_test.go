package env

import (
	"github.com/blend/go-sdk/assert"
	"testing"
)

func TestDelimitedString(t *testing.T) {
	assert := assert.New(t)
	testVars := make(Vars)
	testVars["var_1"] = "val_1"
	res := testVars.DelimitedString(SemicolonDelimiter)
	groundTruth := "var_1=val_1;"

	assert.Equal(groundTruth, res)

	// Now try with multiple key-val pairs
	testVars["var_2"] = "val_2"
	groundTruths := []string{"var_1=val_1;var_2=val_2", "var_2=val_2;var_1=val_1"}

	if !matchOne(assert, t, res, groundTruths...) {
	}
}

func TestParseGoodInputs(t *testing.T) {
	assert := assert.New(t)

	// Single valid key-val pair
	input := "var_1=val_1;"
	res, err := Parse(input, SemicolonDelimiter)
	assert.Equal(nil, err)

	groundTruth := make(Vars)
	groundTruth["var_1"] = "val_1"
	assert.Equal(groundTruth, res)

	// Single valid key-val pair with no trailing delimiter
	input = "var_1=val_1"
	res, err = Parse(input, SemicolonDelimiter)
	assert.Equal(nil, err)

	groundTruth = make(Vars)
	groundTruth["var_1"] = "val_1"
	assert.Equal(groundTruth, res)

	// Two valid key-val pairs
	input = "var_1=val_1;var_2=val_2;"
	res, err = Parse(input, SemicolonDelimiter)
	assert.Equal(nil, err)

	groundTruth = make(Vars)
	groundTruth["var_1"] = "val_1"
	groundTruth["var_2"] = "val_2"
	assert.Equal(groundTruth, res)

	// Two valid key-val pairs with arbitrary whitespace
	input = " var_1   = val_1 ;  var_2 =    val_2; "
	res, err = Parse(input, SemicolonDelimiter)
	assert.Equal(nil, err)
	groundTruth = make(Vars)
	groundTruth["var_1"] = "val_1"
	groundTruth["var_2"] = "val_2"
	assert.Equal(groundTruth, res)

	// Two valid key-val pairs with a quoted string and arbitrary whitespace
	input = "var_1 = val_1; var_2 = \" val_2 \";"
	res, err = Parse(input, SemicolonDelimiter)
	assert.Equal(nil, err)

	groundTruth = make(Vars)
	groundTruth["var_1"] = "val_1"
	groundTruth["var_2"] = " val_2 "
	assert.Equal(groundTruth, res)

	// Two valid key-val pairs with a quoted string and arbitrary whitespace
	// and no trailing separator
	input = `var_1 = val_1; var_2 = " val_2 "`
	res, err = Parse(input, SemicolonDelimiter)
	assert.Equal(nil, err)

	groundTruth = make(Vars)
	groundTruth["var_1"] = "val_1"
	groundTruth["var_2"] = " val_2 "
	assert.Equal(groundTruth, res)

	// Two valid key-val pairs with an escaped quote and arbitrary whitespace
	// and no trailing separator
	input = `var_1 = val_1; var_2 = \" val_2  `
	res, err = Parse(input, SemicolonDelimiter)
	assert.Equal(nil, err)

	groundTruth = make(Vars)
	groundTruth["var_1"] = "val_1"
	groundTruth["var_2"] = `"val_2`
	assert.Equal(groundTruth, res)

	// Two valid key-val pairs with an escaped quote and arbitrary whitespace
	// and no trailing separator
	input = `var_1 = \=val_1; var_2 = \" val_2 \;  `
	res, err = Parse(input, SemicolonDelimiter)
	assert.Equal(nil, err)

	groundTruth = make(Vars)
	groundTruth["var_1"] = "=val_1"
	groundTruth["var_2"] = `"val_2;`
	assert.Equal(groundTruth, res)
}

// matchOne checks to see if the input string is an exact match of a number of
// candidate ground truths
func matchOne(a *assert.Assertions, t *testing.T, input string, groundTruths ...string) bool {
	for _, groundTruth := range groundTruths {
		if groundTruth == input {
			return true
		}
	}
	return false
}
