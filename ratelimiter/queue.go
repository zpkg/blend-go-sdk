/*

Copyright (c) 2022 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package ratelimiter

import (
	"time"

	"github.com/zpkg/blend-go-sdk/collections"
)

// NewQueue returns a new queue based rate limiter.
func NewQueue(numberOfActions int, quantum time.Duration) *Queue {
	return &Queue{
		NumberOfActions: numberOfActions,
		Quantum:         quantum,
		Limits:          map[string]collections.Queue{},
		Now:             func() time.Time { return time.Now().UTC() },
	}
}

// Queue is a simple implementation of a rate checker.
type Queue struct {
	NumberOfActions int
	Quantum         time.Duration
	Limits          map[string]collections.Queue
	Now             func() time.Time
}

// Check returns true if it has been called NumberOfActions times or more in Quantum or smaller duration.
func (q *Queue) Check(id string) bool {
	queue, hasQueue := q.Limits[id]
	if !hasQueue {
		queue = collections.NewRingBufferWithCapacity(q.NumberOfActions)
		q.Limits[id] = queue
	}

	now := q.Now()
	queue.Enqueue(now)
	if queue.Len() < q.NumberOfActions {
		return false
	}

	oldest := queue.Dequeue().(time.Time)
	return now.Sub(oldest) < q.Quantum
}
