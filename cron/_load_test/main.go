package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/blend/go-sdk/cron"
)

const (
	// N is the number of jobs to load.
	N = 1000

	// Q is the total simulation time.
	Q = 30 * time.Second

	// JobRunEvery is the job interval.
	JobRunEvery = 5 * time.Second

	// JobTimeout is the timeout for the jobs.
	JobTimeout = 3 * time.Second

	// JobShortRunTime is the short run time.
	JobShortRunTime = 2 * time.Second

	// JobLongRunTime is the long run time (will induce a timeout.)
	JobLongRunTime = 8 * time.Second
)

var startedCount = new(cron.AtomicCounter)
var completeCount = new(cron.AtomicCounter)
var timeoutCount = new(cron.AtomicCounter)

func ts() string {
	return time.Now().String()
}

type loadTestJob struct {
	id      int
	running bool
	started time.Time
}

func (j *loadTestJob) Timeout() time.Duration {
	return JobTimeout
}

func (j *loadTestJob) Name() string {
	return fmt.Sprintf("loadTestJob_%d", j.id)
}

func (j *loadTestJob) Execute(ctx context.Context) error {
	started := time.Now()
	startedCount.Increment()
	j.running = true

	var runFor time.Duration
	var randValue = rand.Float64()
	if randValue <= 0.5 { // 50% split between short vs. long.
		runFor = JobShortRunTime
	} else {
		runFor = JobLongRunTime
	}

	for time.Now().Sub(started) < runFor {
		time.Sleep(10 * time.Millisecond)
	}
	j.running = false
	completeCount.Increment()
	return nil
}

func (j *loadTestJob) OnCancellation() {
	timeoutCount.Increment()
	j.running = false
}

func (j *loadTestJob) Status() string {
	if j.running {
		return "Request in progress."
	}
	return "Request idle."
}

func (j *loadTestJob) Schedule() cron.Schedule {
	return cron.Every(JobRunEvery)
}

func main() {
	defer func() {
		cron.Default().Stop()
	}()

	if JobLongRunTime < JobTimeout {
		fmt.Printf("Long Run Time: %v is less than the Time Out: %v\n", JobTimeout, JobLongRunTime)
		fmt.Printf("This will cause the Completed vs. Timed Out counts to be wrong.\n")
		os.Exit(1)
	}

	for x := 0; x < N; x++ {
		cron.Default().LoadJob(&loadTestJob{id: x})
	}
	fmt.Printf("Loaded %d Job Instances.\n\n", N)
	cron.Default().Start()

	time.Sleep(Q)

	//given 30 seconds total
	// and running every 5 seconds
	// we expect each job to run 5 times (ish)

	expectedStarted := N * ((int64(Q) / int64(JobRunEvery)) - 1)
	expectedCompleted := expectedStarted >> 1
	expectedTimedOut := expectedStarted >> 1

	fmt.Printf("\nExpected Jobs Started:   %d\n", expectedStarted)
	fmt.Printf("Actual Jobs Started:     %d\n\n", startedCount.Get())

	fmt.Printf("Expected Jobs Completed: %d\n", expectedCompleted)
	fmt.Printf("Actual Jobs Completed:   %d\n\n", completeCount.Get())

	fmt.Printf("Expected Jobs Timed Out: %d\n", expectedTimedOut)
	fmt.Printf("Actual Jobs Timed Out:   %d\n", timeoutCount.Get())
}
