package migration

import "github.com/blend/go-sdk/logger"

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

func New(groups ...*Groups) *Suite {
	return &Migration{
		groups: groups,
	}
}

// Suite is a migration suite.
type Suite struct {
	log *logger.Logger

	groups []*Group

	applied int
	skipped int
	failed  int
	total   int
}

// Logger returns the underlying logger.
func (s *Suite) Logger() *logger.Logger {
	return s.log
}

// Apply applies the suite.
func (s *Suite) Apply(c *db.Connection) (err error) {
	defer s.writeStats()
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", err)
		}
	}()
}



// Applyf active actions to the log.
func (s *Suite) applyf(step *Step, body string, args ...interface{}) {
	if s == nil {
		return
	}
	s.applied = c.applied + 1
	s.total = c.total + 1
	s.write(step, StatApplied, fmt.Sprintf(body, args...))
}

// Skipf passive actions to the log.
func (s *Suite) skipf(m Migration, body string, args ...interface{}) {
	if s == nil {
		return
	}
	s.skipped = s.skipped + 1
	s.total = s.total + 1
	s.write(step, StatSkipped, fmt.Sprintf(body, args...))
}


// Errorf writes errors to the log.
func (s *Suite) errorf(step *Step, body string, args ...interface{}) error {
	if s == nil {
		return 
	}
	s.failed = c.failed + 1
	s.total = c.total + 1
	s.write(m, StatFailed, fmt.Sprintf(body, args...))
}

func (s *Suite) write(step *Step, label, body string) {
	if c == nil {
		return
	}
	if c.output == nil {
		return
	}
	s.log.SyncTrigger(NewEvent(result, body, c.labels(m)...))
}

func (c *Collector) labels(m Migration) []string {
	if c == nil {
		return nil
	}

	var labels []string
	if len(m.Label()) > 0 {
		labels = append(labels, m.Label())
	}

	cursor := m.Parent()
	for cursor != nil {
		if len(cursor.Label()) > 0 {
			labels = append([]string{cursor.Label()}, labels...)
		}
		cursor = cursor.Parent()
	}
	return labels
}