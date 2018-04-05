package workqueue

import (
	"sync"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestWorkerStart(t *testing.T) {
	assert := assert.New(t)
	q := New()
	w := NewWorker(1, q, 1)
	assert.Equal(1, w.ID)
	assert.NotNil(w.Work)

	w.Start()

	wg := sync.WaitGroup{}
	wg.Add(1)
	w.Work <- &Entry{
		Action: func(state ...interface{}) error {
			wg.Done()
			assert.NotEmpty(state)
			assert.Equal("hello", state[0])
			return nil
		},
		Args: []interface{}{"hello"},
	}
	wg.Wait()
}
