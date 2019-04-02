package profanity

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestProfanityRulesFromPath(t *testing.T) {
	assert := assert.New(t)

	profanity := &Profanity{
		Config: &Config{},
	}

	rules, err := profanity.RulesFromPath("../PROFANITY_RULES")
	assert.Nil(err)
	assert.NotEmpty(rules)
}

func TestProfanityReadRules(t *testing.T) {
	assert := assert.New(t)

	profanity := &Profanity{
		Config: &Config{
			RulesFile: "PROFANITY_RULES",
		},
	}

	rules, err := profanity.ReadRules("../")
	assert.Nil(err)
	assert.NotEmpty(rules)
}
