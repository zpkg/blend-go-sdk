package names

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestNames(t *testing.T) {
	assert := assert.New(t)

	names := map[string]*Name{}
	names["John Doe"] = &Name{"", "John", "", "Doe", ""}
	names["Mr Anthony R Von Fange III"] = &Name{"Mr.", "Anthony", "R", "Von Fange", "III"}
	names["Sara Ann Fraser"] = &Name{"", "Sara Ann", "", "Fraser", ""}
	names["Adam"] = &Name{"", "Adam", "", "", ""}
	names["Jonathan Smith"] = &Name{"", "Jonathan", "", "Smith", ""}
	names["Anthony R Von Fange III"] = &Name{"", "Anthony", "R", "Von Fange", "III"}
	names["Anthony Von Fange III"] = &Name{"", "Anthony", "", "Von Fange", "III"}
	names["Mr John Doe"] = &Name{"Mr.", "John", "", "Doe", ""}
	names["Justin White Phd"] = &Name{"", "Justin", "", "White", "PhD"}
	names["Mark P Williams"] = &Name{"", "Mark", "P", "Williams", ""}
	names["Aaron bin Omar"] = &Name{"", "Aaron", "", "bin Omar", ""}
	names["Aaron ibn Omar"] = &Name{"", "Aaron", "", "ibn Omar", ""}

	for rawName, expectedResult := range names {
		result := Parse(rawName)
		assert.Equal(expectedResult, result)
	}
}
