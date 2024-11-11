/*

Copyright (c) 2022 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package redisstats

import (
	"context"

	"github.com/zpkg/blend-go-sdk/logger"
	"github.com/zpkg/blend-go-sdk/redis"
	"github.com/zpkg/blend-go-sdk/stats"
	"github.com/zpkg/blend-go-sdk/timeutil"
)

// AddListeners adds db listeners.
func AddListeners(log logger.Listenable, collector stats.Collector, opts ...stats.AddListenerOption) {
	if log == nil || collector == nil {
		return
	}

	options := stats.NewAddListenerOptions(opts...)

	log.Listen(redis.Flag, stats.ListenerNameStats, redis.NewEventListener(func(ctx context.Context, e redis.Event) {
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
			tags = append(tags, stats.TagError)
		}

		tags = append(tags, options.GetLoggerLabelsAsTags(ctx)...)

		_ = collector.Increment(MetricName, tags...)
		_ = collector.Gauge(MetricNameElapsedLast, timeutil.Milliseconds(e.Elapsed), tags...)
		_ = collector.Histogram(MetricNameElapsed, timeutil.Milliseconds(e.Elapsed), tags...)
	}))
}
