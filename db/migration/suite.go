package migration

import (
	"context"
	"fmt"

	"github.com/blend/go-sdk/db"
	"github.com/blend/go-sdk/exception"
	"github.com/blend/go-sdk/logger"
)

const (
	// StatApplied is a stat name.
	StatApplied = "applied"
	// StatFailed is a stat name.
	StatFailed = "failed"
	// StatSkipped is a stat name.
	StatSkipped = "skipped"
	// StatTotal is a stat name.
	StatTotal = "total"
)

// New returns a new suite of groups.
func New(groups ...GroupedActions) *Suite {
	return &Suite{
		ctx:    context.Background(),
		groups: groups,
	}
}

// Suite is a migration suite.
type Suite struct {
	ctx    context.Context
	log    logger.Log
	groups []GroupedActions

	applied int
	skipped int
	failed  int
	total   int
}

// WithContext sets the suite context.
func (s *Suite) WithContext(ctx context.Context) *Suite {
	s.ctx = ctx
	return s
}

// Context returns the suite context.
func (s *Suite) Context() context.Context {
	return s.ctx
}

// WithLogger sets the suite logger.
func (s *Suite) WithLogger(log logger.Log) *Suite {
	s.log = log
	return s
}

// Logger returns the underlying logger.
func (s *Suite) Logger() logger.Log {
	return s.log
}

// WithGroups adds groups to the suite and returns the suite.
func (s *Suite) WithGroups(groups ...GroupedActions) *Suite {
	s.groups = append(s.groups, groups...)
	return s
}

// Apply applies the suite.
func (s *Suite) Apply(c *db.Connection) (err error) {
	defer s.WriteStats()
	defer func() {
		if r := recover(); r != nil {
			err = exception.New(r)
		}
	}()

	for _, group := range s.groups {
		if err = group.Action(WithSuite(s.Context(), s), c); err != nil {
			return
		}
	}
	return
}

// Applyf writes an applied step message.
func (s *Suite) Applyf(ctx context.Context, format string, args ...interface{}) {
	s.applied = s.applied + 1
	s.total = s.total + 1
	s.Write(ctx, StatApplied, fmt.Sprintf(format, args...))
}

// Skipf skips a given step.
func (s *Suite) Skipf(ctx context.Context, format string, args ...interface{}) {
	s.skipped = s.skipped + 1
	s.total = s.total + 1
	s.Write(ctx, StatSkipped, fmt.Sprintf(format, args...))
}

// Errorf writes an error for a given step.
func (s *Suite) Errorf(ctx context.Context, format string, args ...interface{}) {
	s.failed = s.failed + 1
	s.total = s.total + 1
	s.Write(ctx, StatFailed, fmt.Sprintf(format, args...))
}

// Error
func (s *Suite) Error(ctx context.Context, err error) error {
	s.failed = s.failed + 1
	s.total = s.total + 1
	s.Write(ctx, StatFailed, fmt.Sprintf("%v", err))
	return err
}

func (s *Suite) Write(ctx context.Context, result, body string) {
	logger.MaybeSyncTrigger(s.log, NewEvent(result, body, GetContextLabels(ctx)...))
}

// WriteStats writes the stats if a logger is configured.
func (s *Suite) WriteStats() {
	if s.log != nil {
		s.log.SyncTrigger(NewStatsEvent(s.applied, s.skipped, s.failed, s.total))
	}
}
