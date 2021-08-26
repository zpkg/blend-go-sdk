/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package redisstats

import (
	"context"

	"github.com/blend/go-sdk/ex"
	"github.com/blend/go-sdk/logger"
	"github.com/blend/go-sdk/redis"
	"github.com/blend/go-sdk/stats"
	"github.com/blend/go-sdk/timeutil"
)

// AddListeners adds db listeners.
func AddListeners(log logger.Listenable, collector stats.Collector) {
	if log == nil || collector == nil {
		return
	}

	log.Listen(redis.Flag, stats.ListenerNameStats, redis.NewEventListener(func(_ context.Context, e redis.Event) {
		var tags []string
		if len(e.Network) > 0 {
			tags = append(tags, stats.Tag(TagNetwork, e.Network))
		}
		if len(e.Addr) > 0 {
			tags = append(tags, stats.Tag(TagAddr, e.Addr))
		}
		if len(e.DB) > 0 {
			tags = append(tags, stats.Tag(TagDB, e.DB))
		}
		if len(e.Op) > 0 {
			tags = append(tags, stats.Tag(TagOp, e.Op))
		}
		if e.Err != nil {
			if ex := ex.As(e.Err); ex != nil && ex.Class != nil {
				tags = append(tags, stats.Tag(stats.TagClass, ex.Class.Error()))
			}
			tags = append(tags, stats.TagError)
		}
		_ = collector.Increment(MetricName, tags...)
		_ = collector.Gauge(MetricNameElapsedLast, timeutil.Milliseconds(e.Elapsed), tags...)
		_ = collector.Histogram(MetricNameElapsed, timeutil.Milliseconds(e.Elapsed), tags...)
	}))
}
