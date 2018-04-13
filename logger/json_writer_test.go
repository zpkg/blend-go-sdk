package logger

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestJSONWriter(t *testing.T) {
	assert := assert.New(t)

	output := bytes.NewBuffer(nil)
	jw := NewJSONWriter(output).WithPretty(false)
	assert.False(jw.Pretty())
	assert.Nil(jw.Write(Messagef(Info, "test")))

	var verify JSONObj
	assert.Nil(json.Unmarshal(output.Bytes(), &verify))

	assert.Equal(Info, verify[JSONFieldFlag])
	assert.Equal("test", verify["message"])
}

func TestJSONWriterPretty(t *testing.T) {
	assert := assert.New(t)

	output := bytes.NewBuffer(nil)
	jw := NewJSONWriter(output).WithPretty(true).WithIncludeTimestamp(false)
	assert.True(jw.Pretty())
	assert.False(jw.IncludeTimestamp())
	assert.Nil(jw.Write(Messagef(Info, "test")))

	var verify JSONObj
	assert.Nil(json.Unmarshal(output.Bytes(), &verify))

	assert.Equal(Info, verify[JSONFieldFlag])
	assert.Equal("test", verify["message"])
}
