/*

Copyright (c) 2022 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package cron

import (
	"time"
)

// Schedule is a type that provides a next runtime after a given previous runtime.
type Schedule interface {
	// GetNextRuntime should return the next runtime after a given previous runtime. If `after` is time.Time{} it should be assumed
	// the job hasn't run yet. If time.Time{} is returned by the schedule it is inferred that the job should not run again.
	Next(time.Time) time.Time
}
