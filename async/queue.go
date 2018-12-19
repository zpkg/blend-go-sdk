package async

const (
	// DefaultQueueWorkerMaxWork is the maximum number of work items before queueing blocks.
	DefaultQueueWorkerMaxWork = 1 << 10
)

// NewQueue returns a new queue worker.
func NewQueue(action func(interface{}) error) *QueueWorker {
	return &QueueWorker{
		action:  action,
		latch:   &Latch{},
		maxWork: DefaultQueueWorkerMaxWork,
	}
}

// QueueWorker is a worker that is pushed work over a channel.
type QueueWorker struct {
	action  func(interface{}) error
	latch   *Latch
	errors  chan error
	work    chan interface{}
	maxWork int
}

// WithMaxWork sets the worker max work.
// If set to zero, the worker will block on new work until the last item has processed.
// If set to > 0, the worker will block once the queue reaches the max work length.
// MaxWork will allocate one of each work item to memory.
func (qw *QueueWorker) WithMaxWork(maxWork int) *QueueWorker {
	qw.maxWork = maxWork
	return qw
}

// MaxWork returns the maximum work.
func (qw *QueueWorker) MaxWork() int {
	return qw.maxWork
}

// Latch returns the worker latch.
func (qw *QueueWorker) Latch() *Latch {
	return qw.latch
}

// WithErrors returns the error channel.
func (qw *QueueWorker) WithErrors(errors chan error) *QueueWorker {
	qw.errors = errors
	return qw
}

// Errors returns a channel to read action errors from.
func (qw *QueueWorker) Errors() chan error {
	return qw.errors
}

// Enqueue adds an item to the work queue.
func (qw *QueueWorker) Enqueue(obj interface{}) {
	if qw.work == nil {
		return
	}
	qw.work <- obj
}

// Start starts the worker.
func (qw *QueueWorker) Start() {
	qw.latch.Starting()
	if qw.maxWork > 0 {
		qw.work = make(chan interface{}, qw.maxWork)
	} else {
		qw.work = make(chan interface{})
	}

	go func() {
		qw.latch.Started()
		var err error
		var workItem interface{}
		for {
			select {
			case workItem = <-qw.work:
				err = qw.action(workItem)
				if err != nil && qw.errors != nil {
					qw.errors <- err
				}
			case <-qw.latch.NotifyStopping():
				qw.latch.Stopped()
				return
			}
		}
	}()
	<-qw.latch.NotifyStarted()
}

// Stop stops the worker.
func (qw *QueueWorker) Stop() {
	qw.latch.Stopping()
	<-qw.latch.NotifyStopped()
}
