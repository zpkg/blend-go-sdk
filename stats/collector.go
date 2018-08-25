package stats

// Collector is a stats collector.
type Collector interface {
	DefaultTags() []string
	Count(name string, value int64, tags ...string) error
	Increment(name string, tags ...string) error
	Gauge(name string, value float64, tags ...string) error
	Histogram(name string, value float64, tags ...string) error
}
