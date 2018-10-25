package db

import (
	"testing"

	"github.com/lib/pq"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/exception"
)

func TestError(t *testing.T) {
	assert := assert.New(t)

	assert.Nil(Error(nil))

	var err error
	assert.Nil(Error(err))

	err = exception.New("this is only a test")
	assert.True(exception.Is(Error(err), exception.Class("this is only a test")))

	err = &pq.Error{
		Code:    pq.ErrorCode("P0003"),
		Message: "this is only a test",
		Detail:  "this is only a test",
	}
	assert.NotNil(Error(err))
}
