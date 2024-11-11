/*

Copyright (c) 2022 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package dbstats

import (
	"context"
	"io"
	"testing"
	"time"

	"github.com/zpkg/blend-go-sdk/assert"
	"github.com/zpkg/blend-go-sdk/db"
	"github.com/zpkg/blend-go-sdk/logger"
	"github.com/zpkg/blend-go-sdk/stats"
)

func TestAddListeners(t *testing.T) {
	assert := assert.New(t)

	log := logger.None()
	AddListeners(nil, nil)
	assert.False(log.HasListener(db.QueryFlag, stats.ListenerNameStats))
	AddListeners(log, stats.NewMockCollector(32))
	assert.True(log.HasListener(db.QueryFlag, stats.ListenerNameStats))
}

func TestAddListenersStats(t *testing.T) {
	assert := assert.New(t)

	log := logger.All(logger.OptOutput(io.Discard))
	defer log.Close()
	collector := stats.NewMockCollector(32)

	AddListeners(log, collector)

	log.TriggerContext(context.Background(), db.NewQueryEvent("select 'ok!'", time.Second))

	qm := <-collector.Metrics
	assert.Equal(MetricNameDBQuery, qm.Name)
	assert.Equal(1, qm.Count)
	assert.NotEmpty(qm.Tags)

	qm = <-collector.Metrics
	assert.Equal(MetricNameDBQueryElapsedLast, qm.Name)
	assert.Equal(1000, qm.Gauge)
	assert.NotEmpty(qm.Tags)

	qm = <-collector.Metrics
	assert.Equal(MetricNameDBQueryElapsed, qm.Name)
	assert.Equal(1000, qm.Histogram)
	assert.NotEmpty(qm.Tags)
}
