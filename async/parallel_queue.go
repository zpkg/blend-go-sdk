package async

import (
	"runtime"
	"sync"
)

// NewParallelQueue returns a new parallel queue worker.
func NewParallelQueue(action func(interface{}) error) *ParallelQueue {
	return &ParallelQueue{
		latch:      &Latch{},
		work:       make(chan interface{}, DefaultQueueMaxWork),
		action:     action,
		numWorkers: runtime.NumCPU(),
	}
}

// ParallelQueue is a queude with multiple workers..
type ParallelQueue struct {
	sync.Mutex
	latch      *Latch
	numWorkers int
	workers    chan *Queue
	action     func(interface{}) error
	errors     chan error
	work       chan interface{}
	draining   *sync.WaitGroup
}

// WithWorkers sets the number of workers.
// It defaults to `runtime.NumCPU()`
func (pq *ParallelQueue) WithWorkers(numWorkers int) *ParallelQueue {
	pq.numWorkers = numWorkers
	return pq
}

// Workers returns the number of worker route
func (pq *ParallelQueue) Workers() int {
	return pq.numWorkers
}

// WithWork sets the work channel.
func (pq *ParallelQueue) WithWork(work chan interface{}) *ParallelQueue {
	pq.work = work
	return pq
}

// Work returns the work channel.
func (pq *ParallelQueue) Work() chan interface{} {
	return pq.work
}

// Latch returns the worker latch.
func (pq *ParallelQueue) Latch() *Latch {
	return pq.latch
}

// WithErrors sets the error channel.
func (pq *ParallelQueue) WithErrors(errors chan error) *ParallelQueue {
	pq.errors = errors
	return pq
}

// Errors returns a channel to read action errors from.
// You must provide it with `WithErrors`.
func (pq *ParallelQueue) Errors() chan error {
	return pq.errors
}

// Enqueue adds an item to the work queue.
func (pq *ParallelQueue) Enqueue(obj interface{}) {
	pq.Lock()
	defer pq.Unlock()
	println("queueing item", obj)
	pq.work <- obj
}

// StartWorkers starts all workers.
func (pq *ParallelQueue) StartWorkers() {
	for x := 0; x < pq.numWorkers; x++ {
		worker := <-pq.workers
		println("starting worker", x)
		worker.Start()
		pq.workers <- worker
	}
}

// StopWorkers closes all workers.
func (pq *ParallelQueue) StopWorkers() {
	for x := 0; x < pq.numWorkers; x++ {
		worker := <-pq.workers
		println("stopping worker", x)
		worker.Stop()
		pq.workers <- worker
	}
}

// InitializeWorkers initializes the workers.
func (pq *ParallelQueue) InitializeWorkers() {
	println("initializing parallel queue workers")
	pq.workers = make(chan *Queue, pq.numWorkers)
	for x := 0; x < pq.numWorkers; x++ {
		worker := &Queue{
			latch:  &Latch{},
			errors: pq.errors,
			work:   make(chan interface{}),
		}
		worker.action = pq.AndReturn(worker, pq.action)
		pq.workers <- worker
	}
}

// Start starts the worker.
func (pq *ParallelQueue) Start() {
	println("starting parallel queue")
	pq.latch.Starting()
	// if not initialized, initialize workers
	if pq.workers == nil {
		pq.InitializeWorkers()
	}
	pq.StartWorkers()
	go pq.Dispatch()
	<-pq.latch.NotifyStarted()
}

// Dispatch processes work items in a loop.
func (pq *ParallelQueue) Dispatch() {
	println("dispatch started")
	pq.latch.Started()
	var workItem interface{}
	var worker *Queue
	for {
		select {
		case workItem = <-pq.work:
			select {
			case worker = <-pq.workers:
				worker.work <- workItem
			case <-pq.latch.NotifyStopping():
				pq.latch.Stopped()
				return
			}
		case <-pq.latch.NotifyStopping():
			pq.latch.Stopped()
			println("dispatch stopping")
			return
		}
	}
}

// AndReturn creates an action handler that returns a given worker to the worker queue.
// It wraps any action provided to the queue.
func (pq *ParallelQueue) AndReturn(worker *Queue, action func(interface{}) error) func(interface{}) error {
	return func(workItem interface{}) error {
		defer func() {
			if pq.draining != nil {
				println("completing drained work item")
				pq.draining.Done()
			} else {
				println("completing work item")
			}
			pq.workers <- worker
		}()
		return action(workItem)
	}
}

// Drain drains the queue.
func (pq *ParallelQueue) Drain() error {
	println("draining")

	pq.Lock()
	defer pq.Unlock()

	println("stopping workers")
	pq.StopWorkers()

	println("stopping dispatch")
	pq.latch.Stopping()
	<-pq.latch.NotifyStopped()

	println(len(pq.work), "items left")
	pq.draining = &sync.WaitGroup{}
	pq.draining.Add(len(pq.work))

	println("restarting workers")
	pq.StartWorkers()

	println("restarting dispatch loop")
	pq.latch.Starting()
	go pq.Dispatch()
	<-pq.latch.NotifyStarted()

	println("waiting for work to complete")
	pq.draining.Wait()
	pq.draining = nil

	println("restarted, draining complete")
	return nil
}

// Close stops the queue.
// Any work left in the queue will be discarded.
func (pq *ParallelQueue) Close() error {
	pq.latch.Stopping()
	<-pq.latch.NotifyStopped()

	for x := 0; x < pq.numWorkers; x++ {
		worker := <-pq.workers
		if err := worker.Close(); err != nil {
			return err
		}
	}
	return nil
}
