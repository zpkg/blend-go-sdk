package async

import (
	"sync"
)

// NewParallelQueue returns a new parallel queue worker.
func NewParallelQueue(numWorkers int, action func(interface{}) error) *ParallelQueue {
	return &ParallelQueue{
		action:     action,
		numWorkers: numWorkers,
		latch:      &Latch{},
		workers:    make(chan *Queue, numWorkers),
		work:       make(chan interface{}, DefaultQueueMaxWork),
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
	if pq.draining != nil {
		pq.draining.Wait()
	}
	pq.work <- obj
	pq.Unlock()
}

// Start starts the worker.
func (pq *ParallelQueue) Start() {
	pq.latch.Starting()
	for x := 0; x < pq.numWorkers; x++ {
		worker := &Queue{
			latch:  &Latch{},
			errors: pq.errors,
			work:   make(chan interface{}),
		}
		worker.action = pq.AndReturn(worker, pq.action)
		worker.Start()
		pq.workers <- worker
	}
	go pq.Dispatch()
	<-pq.latch.NotifyStarted()
}

// Dispatch processes work items in a loop.
func (pq *ParallelQueue) Dispatch() {
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
				pq.draining.Done()
			}
			pq.workers <- worker
		}()
		return action(workItem)
	}
}

// Drain drains the queue.
func (pq *ParallelQueue) Drain() error {
	pq.Lock()
	defer pq.Unlock()

	pq.draining = &sync.WaitGroup{}
	pq.draining.Add(len(pq.work))
	pq.draining.Wait()
	pq.draining = nil

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
