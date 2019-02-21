package secrets

import (
	"context"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestMockClient(t *testing.T) {
	assert := assert.New(t)
	todo := context.TODO()

	client := NewMockClient()

	err := client.Put(todo, "testo", map[string]string{"key_123": "value_xyz"})
	assert.Nil(err)

	vals, err := client.Get(todo, "testo")
	assert.Nil(err)
	assert.Equal("value_xyz", vals["key_123"])

	_, err = client.Get(todo, "fake_test")
	assert.NotNil(err)

	err = client.Delete(todo, "another_fake")
	assert.NotNil(err)

	err = client.Delete(todo, "testo")
	assert.Nil(err)

	_, err = client.Get(todo, "testo")
	assert.NotNil(err)
}
