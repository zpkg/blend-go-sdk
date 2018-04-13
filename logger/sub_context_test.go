package logger

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestLoggerSubContext(t *testing.T) {
	assert := assert.New(t)

	l := New().WithHeading("test-logger")
	sc := l.SubContext("sub-context")
	assert.NotNil(sc.Logger())
	assert.Equal([]string{"test-logger", "sub-context"}, sc.Headings())
}
