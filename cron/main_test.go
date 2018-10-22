package cron

import (
	"context"
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
)

// TestMain is the testing entrypoint.
func TestMain(m *testing.M) {
	assert.Main(m)
}

const (
	runAtJobName = "runAt"
)

type runAtJob struct {
	RunAt       time.Time
	RunDelegate func(ctx context.Context) error
}

type runAt time.Time

func (ra runAt) GetNextRunTime(after *time.Time) *time.Time {
	typed := time.Time(ra)
	return &typed
}

func (raj *runAtJob) Name() string {
	return "runAt"
}

func (raj *runAtJob) Schedule() Schedule {
	return runAt(raj.RunAt)
}

func (raj *runAtJob) Execute(ctx context.Context) error {
	return raj.RunDelegate(ctx)
}

var (
	_ OnCancellationReceiver = (*testJobWithTimeout)(nil)
)

type testJobWithTimeout struct {
	RunAt                time.Time
	TimeoutDuration      time.Duration
	RunDelegate          func(ctx context.Context) error
	CancellationDelegate func()
}

func (tj *testJobWithTimeout) Name() string {
	return "testJobWithTimeout"
}

func (tj *testJobWithTimeout) Timeout() time.Duration {
	return tj.TimeoutDuration
}

func (tj *testJobWithTimeout) Schedule() Schedule {
	return Immediately()
}

func (tj *testJobWithTimeout) Execute(ctx context.Context) error {
	return tj.RunDelegate(ctx)
}

func (tj *testJobWithTimeout) OnCancellation(ji *JobInvocation) {
	tj.CancellationDelegate()
}

type testJobInterval struct {
	RunEvery    time.Duration
	RunDelegate func(ctx context.Context) error
}

func (tj *testJobInterval) Name() string {
	return "testJobInterval"
}

func (tj *testJobInterval) Schedule() Schedule {
	return Every(tj.RunEvery)
}

func (tj *testJobInterval) Execute(ctx context.Context) error {
	return tj.RunDelegate(ctx)
}

type testWithEnabled struct {
	isEnabled bool
	action    func()
}

func (twe testWithEnabled) Name() string {
	return "testWithEnabled"
}

func (twe testWithEnabled) Schedule() Schedule {
	return nil
}

func (twe testWithEnabled) Enabled() bool {
	return twe.isEnabled
}

func (twe testWithEnabled) Execute(ctx context.Context) error {
	twe.action()
	return nil
}

type mockTracer struct {
	OnStart  func(*JobInvocation)
	OnFinish func(*JobInvocation)
}

func (mt mockTracer) Start(ctx context.Context, ji *JobInvocation) (context.Context, TraceFinisher) {
	if mt.OnStart != nil {
		mt.OnStart(ji)
	}
	return ctx, &mockTraceFinisher{Parent: &mt}
}

type mockTraceFinisher struct {
	Parent *mockTracer
}

func (mtf mockTraceFinisher) Finish(ctx context.Context, ji *JobInvocation) {
	if mtf.Parent != nil && mtf.Parent.OnFinish != nil {
		mtf.Parent.OnFinish(ji)
	}
}

var (
	_ OnStartReceiver    = (*brokenFixedTest)(nil)
	_ OnCompleteReceiver = (*brokenFixedTest)(nil)
	_ OnFailureReceiver  = (*brokenFixedTest)(nil)
	_ OnBrokenReceiver   = (*brokenFixedTest)(nil)
	_ OnFixedReceiver    = (*brokenFixedTest)(nil)
)

func newBrokenFixedTest(action func(context.Context) error) *brokenFixedTest {
	return &brokenFixedTest{
		BrokenSignal: make(chan struct{}),
		FixedSignal:  make(chan struct{}),
		Action:       action,
	}
}

type brokenFixedTest struct {
	Starts       int
	Completes    int
	Failures     int
	BrokenSignal chan struct{}
	FixedSignal  chan struct{}
	Action       func(context.Context) error
}

func (job brokenFixedTest) Name() string { return "broken-fixed" }

func (job brokenFixedTest) Execute(ctx context.Context) error {
	return job.Action(ctx)
}

func (job brokenFixedTest) Schedule() Schedule { return nil }

func (job *brokenFixedTest) OnStart(t *JobInvocation) {
	job.Starts++
}

func (job *brokenFixedTest) OnFailure(t *JobInvocation) {
	job.Failures++
}

func (job *brokenFixedTest) OnComplete(t *JobInvocation) {
	job.Completes++
}

func (job *brokenFixedTest) OnBroken(t *JobInvocation) {
	close(job.BrokenSignal)
}

func (job *brokenFixedTest) OnFixed(t *JobInvocation) {
	close(job.FixedSignal)
}

type loadJobTestMinimum struct{}

func (job loadJobTestMinimum) Name() string                    { return "load-job-test-minimum" }
func (job loadJobTestMinimum) Execute(_ context.Context) error { return nil }
