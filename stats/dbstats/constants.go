/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package dbstats

import "github.com/blend/go-sdk/db"

// Metric and tag names etc.
const (
	MetricNameDBQuery        string = string(db.QueryFlag)
	MetricNameDBQueryElapsed string = MetricNameDBQuery + ".elapsed"

	TagQuery    string = "query"
	TagEngine   string = "engine"
	TagDatabase string = "database"
)
