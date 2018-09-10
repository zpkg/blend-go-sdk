package async

import (
	"fmt"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestRunToError(t *testing.T) {
	assert := assert.New(t)

	fn := func() error {
		return nil
	}
	assert.Nil(RunToError(fn, fn))

	fn = func() error {
		return fmt.Errorf("ERROR")
	}
	assert.NotNil(RunToError(fn, fn))
}
