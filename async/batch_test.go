package async

import (
	"fmt"
	"strconv"
	"sync/atomic"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestBatch(t *testing.T) {
	assert := assert.New(t)

	var work []interface{}
	for x := 0; x < 32; x++ {
		work = append(work, "hello"+strconv.Itoa(x))
	}

	var processed int32
	errors := make(chan error, 32)
	b := NewBatch(func(v interface{}) error {
		atomic.AddInt32(&processed, 1)
		return fmt.Errorf("this is only a test")
	}, work...).WithErrors(errors)
	b.Process()

	assert.Equal(32, processed)
	assert.Equal(32, len(errors))
}
