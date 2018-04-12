package cron

import (
	"sync"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestAtomicFlag(t *testing.T) {
	assert := assert.New(t)

	zero := &AtomicFlag{}
	assert.False(zero.Get())
	shouldBeTrue := &AtomicFlag{}
	shouldBeTrue.Set(true)
	assert.True(shouldBeTrue.Get())
	shouldBeFalse := &AtomicFlag{}
	shouldBeFalse.Set(false)
	assert.False(shouldBeFalse.Get())
}

func TestAtomicCounter(t *testing.T) {
	assert := assert.New(t)

	ac := &AtomicCounter{}

	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		for x := 0; x < 100; x++ {
			ac.Increment()
		}
	}()
	go func() {
		defer wg.Done()
		for x := 0; x < 100; x++ {
			ac.Increment()
		}
	}()

	wg.Wait()
	assert.Equal(200, ac.Get())

	wg.Add(2)
	go func() {
		defer wg.Done()
		for x := 0; x < 100; x++ {
			ac.Decrement()
		}
	}()
	go func() {
		defer wg.Done()
		for x := 0; x < 100; x++ {
			ac.Decrement()
		}
	}()
	wg.Wait()
	assert.Zero(ac.Get())
}
