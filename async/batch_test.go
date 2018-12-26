package async

import (
	"context"
	"fmt"
	"strconv"
	"sync/atomic"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestBatch(t *testing.T) {
	assert := assert.New(t)

	var processed int32
	errors := make(chan error, 32)
	b := NewBatch(func(_ context.Context, v interface{}) error {
		atomic.AddInt32(&processed, 1)
		return fmt.Errorf("this is only a test")
	}).WithErrors(errors).WithWork(make(chan interface{}, 32))
	for x := 0; x < 32; x++ {
		b.Add("hello" + strconv.Itoa(x))
	}
	b.ProcessContext(context.Background())

	assert.Equal(32, processed)
	assert.Equal(32, len(errors))
}
