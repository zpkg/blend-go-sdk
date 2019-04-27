package slant

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestPrint(t *testing.T) {
	assert := assert.New(t)

	expected := ` _       _____    ____  ____  _______   __
| |     / /   |  / __ \/ __ \/ ____/ | / /
| | /| / / /| | / /_/ / / / / __/ /  |/ /
| |/ |/ / ___ |/ _, _/ /_/ / /___/ /|  /
|__/|__/_/  |_/_/ |_/_____/_____/_/ |_/`

	output, err := PrintString("WARDEN")
	assert.Nil(err)
	assert.Equal(expected, output)
}
