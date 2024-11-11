/*

Copyright (c) 2022 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package cron

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/zpkg/blend-go-sdk/logger"
	"github.com/zpkg/blend-go-sdk/timeutil"
)

// these are compile time assertions
var (
	_ logger.Event        = (*Event)(nil)
	_ logger.TextWritable = (*Event)(nil)
	_ logger.JSONWritable = (*Event)(nil)
)

// NewEventListener returns a new event listener.
func NewEventListener(listener func(context.Context, Event)) logger.Listener {
	return func(ctx context.Context, e logger.Event) {
		if typed, isTyped := e.(Event); isTyped {
			listener(ctx, typed)
		}
	}
}

// NewEvent creates a new event with a given set of optional options.
func NewEvent(flag, jobName string, options ...EventOption) Event {
	e := Event{
		Flag:    flag,
		JobName: jobName,
	}

	for _, option := range options {
		option(&e)
	}
	return e
}

// EventOption is an option for an Event.
type EventOption func(*Event)

// OptEventJobInvocation sets a field.
func OptEventJobInvocation(jobInvocation string) EventOption {
	return func(e *Event) { e.JobInvocation = jobInvocation }
}

// OptEventErr sets a field.
func OptEventErr(err error) EventOption {
	return func(e *Event) { e.Err = err }
}

// OptEventElapsed sets a field.
func OptEventElapsed(elapsed time.Duration) EventOption {
	return func(e *Event) { e.Elapsed = elapsed }
}

// Event is an event.
type Event struct {
	Flag          string
	JobName       string
	JobInvocation string
	Err           error
	Elapsed       time.Duration
}

// GetFlag implements logger.Event.
func (e Event) GetFlag() string { return e.Flag }

// Complete returns if the event completed.
func (e Event) Complete() bool {
	return e.Flag == FlagComplete
}

// WriteText implements logger.TextWritable.
func (e Event) WriteText(tf logger.TextFormatter, wr io.Writer) {
	if e.Elapsed > 0 {
		fmt.Fprint(wr, logger.Space)
		fmt.Fprintf(wr, "(%v)", e.Elapsed)
	}
}

// Decompose implements logger.JSONWritable.
func (e Event) Decompose() map[string]interface{} {
	return map[string]interface{}{
		"jobName": e.JobName,
		"err":     e.Err,
		"elapsed": timeutil.Milliseconds(e.Elapsed),
	}
}
