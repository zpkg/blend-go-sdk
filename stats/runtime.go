package stats

import (
	"runtime"
	"time"

	"github.com/blend/go-sdk/logger"
)

// Runtime reports golang vm runtime stats.
func Runtime(log *logger.Logger, collector Collector) {
	if collector == nil {
		return
	}

	go func() {
		var previous, current runtime.MemStats
		runtimeCollect(log, collector, &previous, &current)
		for {
			select {
			case <-time.After(250 * time.Millisecond):
				runtimeCollect(log, collector, &previous, &current)
			}
		}
	}()
}

func runtimeCollect(log *logger.Logger, collector Collector, previous, current *runtime.MemStats) {
	runtime.ReadMemStats(current)

	// these depend on the previous values
	collector.Gauge("go.runtime.mem.num_gc", float64(current.NumGC-previous.NumGC))
	collector.Gauge("go.runtime.mem.num_forced_gc", float64(current.NumForcedGC-previous.NumForcedGC))
	collector.Gauge("go.runtime.mem.pause_total_ns", float64(current.PauseTotalNs-previous.PauseTotalNs))
	collector.Gauge("go.runtime.mem.frees", float64(current.Frees-previous.Frees))
	collector.Gauge("go.runtime.mem.mallocs", float64(current.Mallocs-previous.Mallocs))

	// these are mostly points in time.
	collector.Gauge("go.runtime.num_cpu", float64(runtime.NumCPU()))
	collector.Gauge("go.runtime.num_goroutine", float64(runtime.NumGoroutine()))

	collector.Gauge("go.runtime.mem.alloc", float64(current.Alloc))

	collector.Gauge("go.runtime.mem.gc_sys", float64(current.GCSys))
	collector.Gauge("go.runtime.mem.other_sys", float64(current.OtherSys))

	collector.Gauge("go.runtime.mem.heap_alloc", float64(current.HeapAlloc))
	collector.Gauge("go.runtime.mem.heap_idle", float64(current.HeapIdle))
	collector.Gauge("go.runtime.mem.heap_inuse", float64(current.HeapInuse))
	collector.Gauge("go.runtime.mem.heap_objects", float64(current.HeapObjects))
	collector.Gauge("go.runtime.mem.heap_sys", float64(current.HeapSys))

	collector.Gauge("go.runtime.mem.stack_inuse", float64(current.StackInuse))
	collector.Gauge("go.runtime.mem.stack_sys", float64(current.StackSys))
	collector.Gauge("go.runtime.mem.sys", float64(current.Sys))
	collector.Gauge("go.runtime.mem.total_alloc", float64(current.TotalAlloc))

	// rotate the results ...
	*previous = *current
}
