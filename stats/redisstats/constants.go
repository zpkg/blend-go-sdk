/*

Copyright (c) 2022 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package redisstats

import "github.com/zpkg/blend-go-sdk/redis"

// Metric and tag names etc.
const (
	MetricName            string = string(redis.Flag)
	MetricNameElapsed     string = MetricName + ".elapsed"
	MetricNameElapsedLast string = MetricNameElapsed + ".last"

	TagNetwork string = "network"
	TagAddr    string = "addr"
	TagDB      string = "db"
	TagOp      string = "op"
)
