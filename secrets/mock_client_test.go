package secrets

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestMockClient(t *testing.T) {
	assert := assert.New(t)

	client := NewMockClient()

	err := client.Put("testo", map[string]string{"key_123": "value_xyz"})
	assert.Nil(err)

	vals, err := client.Get("testo")
	assert.Nil(err)
	assert.Equal("value_xyz", vals["key_123"])

	_, err = client.Get("fake_test")
	assert.NotNil(err)

	err = client.Delete("another_fake")
	assert.NotNil(err)

	err = client.Delete("testo")
	assert.Nil(err)

	_, err = client.Get("testo")
	assert.NotNil(err)
}
