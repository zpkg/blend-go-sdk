/*

Copyright (c) 2022 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package cron

import (
	"context"

	"github.com/zpkg/blend-go-sdk/logger"
)

// JobManagerOption is a job manager option.
type JobManagerOption func(*JobManager)

// OptLog sets the job manager logger.
func OptLog(log logger.Log) JobManagerOption {
	return func(jm *JobManager) { jm.Log = log }
}

// OptTracer sets the job manager tracer.
func OptTracer(tracer Tracer) JobManagerOption {
	return func(jm *JobManager) { jm.Tracer = tracer }
}

// OptBaseContext sets the job manager base context.
func OptBaseContext(ctx context.Context) JobManagerOption {
	return func(jm *JobManager) { jm.BaseContext = ctx }
}
