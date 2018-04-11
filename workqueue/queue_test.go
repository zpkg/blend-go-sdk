package workqueue

import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/util"
)

// TestEnqueue tests that things are queued and dispatched in order with a single worker.
func TestEnqueue(t *testing.T) {
	assert := assert.New(t)

	q := NewWithWorkers(1)
	defer q.Close()
	q.Start()

	wg := sync.WaitGroup{}
	wg.Add(3)

	var outputs []string
	var outputsLock sync.Mutex
	q.Enqueue(func(workData ...interface{}) error {
		outputsLock.Lock()
		defer outputsLock.Unlock()

		outputs = append(outputs, workData[0].(string))
		wg.Done()
		return nil
	}, "Hello")

	q.Enqueue(func(workData ...interface{}) error {
		outputsLock.Lock()
		defer outputsLock.Unlock()

		outputs = append(outputs, workData[0].(string))
		wg.Done()
		return nil
	}, "Test")

	q.Enqueue(func(workData ...interface{}) error {
		outputsLock.Lock()
		defer outputsLock.Unlock()

		outputs = append(outputs, workData[0].(string))
		wg.Done()
		return nil
	}, "World")

	wg.Wait()
	assert.Len(outputs, 3)
	assert.Equal("Hello", outputs[0], fmt.Sprintf("%#v", outputs))
	assert.Equal("Test", outputs[1], fmt.Sprintf("%#v", outputs))
	assert.Equal("World", outputs[2], fmt.Sprintf("%#v", outputs))
}

func TestDrain(t *testing.T) {
	assert := assert.New(t)

	q := NewWithWorkers(1)

	wg := sync.WaitGroup{}
	wg.Add(1)

	q.Start()
	q.Enqueue(func(workData ...interface{}) error {
		wg.Done()
		return nil
	})
	wg.Wait()
	q.Close()

	assert.False(q.Running())
}

func TestProcessQueueReqeue(t *testing.T) {
	assert := assert.New(t)
	assert.StartTimeout(1*time.Second, "This should take < 1 second")

	defer func() {
		assert.EndTimeout()
	}()

	q := NewWithWorkers(2)
	defer q.Close()

	q.Start()

	maxErrors := q.maxRetries - 1
	numErrors := 0

	wg := sync.WaitGroup{}
	wg.Add(1)

	var output string
	q.Enqueue(func(workData ...interface{}) error {
		if numErrors < maxErrors {
			numErrors = numErrors + 1
			return errors.New("Requeue")
		}

		output = workData[0].(string)
		wg.Done()
		return nil
	}, "Hello")

	wg.Wait()
	assert.Equal(maxErrors, numErrors)
	assert.Equal("Hello", output)
}

func TestProcessQueuePanics(t *testing.T) {
	assert := assert.New(t)
	q := NewWithWorkers(2)
	defer q.Close()
	q.Start()

	wg := sync.WaitGroup{}
	wg.Add(5)

	var numErrors int32
	q.Enqueue(func(workData ...interface{}) error {
		defer wg.Done()
		atomic.AddInt32(&numErrors, 1)

		if numErrors == 2 {
			panic("this is only a test")
		}

		if numErrors < 5 {
			return fmt.Errorf("new error")
		}

		return nil
	}, "Hello")
	wg.Wait()

	assert.True(q.Running())
	assert.Equal(5, numErrors)

	wg2 := sync.WaitGroup{}
	wg2.Add(1)
	q.Enqueue(func(data ...interface{}) error {
		defer wg2.Done()
		assert.Equal("test", data[0].(string))
		return nil
	}, "test")

	wg2.Wait()
	assert.Equal(0, q.Len())
}

func TestProcessQueueParallel(t *testing.T) {
	assert := assert.New(t)
	assert.StartTimeout(1*time.Second, "This should take < 1 second")
	defer func() {
		assert.EndTimeout()
	}()

	q := NewWithWorkers(2)
	defer q.Close()
	q.Start()

	wg := sync.WaitGroup{}
	wg.Add(4)

	var runCount int32
	q.Enqueue(func(workData ...interface{}) error {
		atomic.AddInt32(&runCount, 1)
		wg.Done()
		return nil
	})

	q.Enqueue(func(workData ...interface{}) error {
		atomic.AddInt32(&runCount, 1)
		wg.Done()
		return nil
	})

	q.Enqueue(func(workData ...interface{}) error {
		atomic.AddInt32(&runCount, 1)
		wg.Done()
		return nil
	})

	q.Enqueue(func(workData ...interface{}) error {
		atomic.AddInt32(&runCount, 1)
		wg.Done()
		return nil
	})

	wg.Wait()
	assert.Equal(4, runCount)
}

func benchEnqueue(workers, items int) time.Duration {
	wg := sync.WaitGroup{}
	wg.Add(items)

	q := NewWithWorkers(workers)
	q.SetMaxWorkItems(items)
	q.Start()

	times := make([]time.Duration, items)
	for x := 0; x < items; x++ {
		q.Enqueue(func(workData ...interface{}) error {
			defer wg.Done()
			times[workData[0].(int)] = time.Since(workData[1].(time.Time))
			return nil
		}, x, time.Now())
	}

	wg.Wait()
	q.Close()
	return util.Math.MeanOfDuration(times)
}

func chanToArray(times chan time.Duration) []time.Duration {
	values := make([]time.Duration, len(times))
	for x := 0; x < len(times); x++ {
		values = append(values, <-times)
	}
	return values
}

func BenchmarkQueueEnqueueSync1024(b *testing.B) {
	averageTimes := make([]time.Duration, b.N)
	for n := 0; n < b.N; n++ {
		averageTimes[n] = benchEnqueue(32, 1024)
	}

	b.Logf("Average time in queue: %v\n", util.Math.MeanOfDuration(averageTimes))
}

func BenchmarkQueueEnqueueSync2048_1(b *testing.B) {
	averageTimes := make([]time.Duration, b.N)
	for n := 0; n < b.N; n++ {
		averageTimes[n] = benchEnqueue(1, 2048)
	}

	b.Logf("Average time in queue: %v\n", util.Math.MeanOfDuration(averageTimes))
}

func BenchmarkQueueEnqueueSync2048_2(b *testing.B) {
	averageTimes := make([]time.Duration, b.N)
	for n := 0; n < b.N; n++ {
		averageTimes[n] = benchEnqueue(2, 2048)
	}

	b.Logf("Average time in queue: %v\n", util.Math.MeanOfDuration(averageTimes))
}

func BenchmarkQueueEnqueueSync2048_4(b *testing.B) {
	averageTimes := make([]time.Duration, b.N)
	for n := 0; n < b.N; n++ {
		averageTimes[n] = benchEnqueue(4, 2048)
	}

	b.Logf("Average time in queue: %v\n", util.Math.MeanOfDuration(averageTimes))
}

func BenchmarkQueueEnqueueSync2048_8(b *testing.B) {
	averageTimes := make([]time.Duration, b.N)
	for n := 0; n < b.N; n++ {
		averageTimes[n] = benchEnqueue(8, 2048)
	}

	b.Logf("Average time in queue: %v\n", util.Math.MeanOfDuration(averageTimes))
}

func BenchmarkQueueEnqueueSync2048_32(b *testing.B) {
	averageTimes := make([]time.Duration, b.N)
	for n := 0; n < b.N; n++ {
		averageTimes[n] = benchEnqueue(32, 2048)
	}

	b.Logf("Average time in queue: %v\n", util.Math.MeanOfDuration(averageTimes))
}

func BenchmarkQueueEnqueueSync4096(b *testing.B) {
	averageTimes := make([]time.Duration, b.N)
	for n := 0; n < b.N; n++ {
		averageTimes[n] = benchEnqueue(32, 4096)
	}

	b.Logf("Average time in queue: %v\n", util.Math.MeanOfDuration(averageTimes))
}

func BenchmarkQueueEnqueueSync8192(b *testing.B) {
	averageTimes := make([]time.Duration, b.N)
	for n := 0; n < b.N; n++ {
		averageTimes[n] = benchEnqueue(32, 8192)
	}

	b.Logf("Average time in queue: %v\n", util.Math.MeanOfDuration(averageTimes))
}

func BenchmarkQueueEnqueueSync16384(b *testing.B) {
	averageTimes := make([]time.Duration, b.N)
	for n := 0; n < b.N; n++ {
		averageTimes[n] = benchEnqueue(32, 16384)
	}

	b.Logf("Average time in queue: %v\n", util.Math.MeanOfDuration(averageTimes))
}

func BenchmarkQueueEnqueueSync32768(b *testing.B) {
	averageTimes := make([]time.Duration, b.N)
	for n := 0; n < b.N; n++ {
		averageTimes[n] = benchEnqueue(64, 32768)
	}

	b.Logf("Average time in queue: %v\n", util.Math.MeanOfDuration(averageTimes))
}
