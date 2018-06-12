package exception

import (
	"encoding/json"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestErrorMarshalJSON(t *testing.T) {
	assert := assert.New(t)

	err := Error("this is only a test")
	contents, marshalErr := json.Marshal(err)
	assert.Nil(marshalErr)
	assert.Equal(`"this is only a test"`, string(contents))
}
