package profanity

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestContainsAny(t *testing.T) {
	assert := assert.New(t)

	ruleFunc := ContainsAny("foo", "bar")

	assert.NotNil(ruleFunc("", []byte(`aaa foo`)))
	assert.NotNil(ruleFunc("", []byte(`foo aaa`)))
	assert.NotNil(ruleFunc("", []byte(`aaa foo aaa`)))

	assert.NotNil(ruleFunc("", []byte(`aaa bar`)))
	assert.NotNil(ruleFunc("", []byte(`bar aaa`)))
	assert.NotNil(ruleFunc("", []byte(`aaa bar aaa`)))

	assert.Nil(ruleFunc("", []byte(``)))
	assert.Nil(ruleFunc("", []byte(`aaa`)))
}

func TestContainsAnyReportsLineNumber(t *testing.T) {
	assert := assert.New(t)

	ruleFunc := ContainsAny("foo")

	file := `111
222
foo
444
555
`

	err := ruleFunc("", []byte(file))
	assert.Contains(err.Error(), "line: 3")
}
