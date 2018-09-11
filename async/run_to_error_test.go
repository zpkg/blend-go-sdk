package async

import (
	"fmt"
	"sync"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestRunToErrorError(t *testing.T) {
	assert := assert.New(t)

	hostedFinished := make(chan struct{})
	finished := make(chan struct{})

	var err error
	go func() {
		defer func() {
			close(finished)
		}()
		err = RunToError(func() error {
			defer func() {
				close(hostedFinished)
			}()
			return fmt.Errorf("an error")
		})
	}()

	<-hostedFinished
	<-finished
	assert.NotNil(err)
	assert.Equal("an error", err.Error())
}

func TestRunToErrorPanic(t *testing.T) {
	assert := assert.New(t)

	wg := sync.WaitGroup{}
	wg.Add(2)

	stop := make(chan struct{})
	finished := make(chan struct{})

	var err error
	go func() {
		defer func() {
			if r := recover(); r != nil {
				err = fmt.Errorf("%v", r)
			}
			close(finished)
		}()

		err = RunToError(func() error {
			defer wg.Done()
			<-stop
			return nil
		}, func() error {
			defer wg.Done()
			<-stop
			panic("only a test")
		})
	}()

	close(stop)
	<-finished
	wg.Wait()
	assert.NotNil(err)
	assert.Equal("only a test", err.Error())
}
