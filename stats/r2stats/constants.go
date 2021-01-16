/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package r2stats

import "github.com/blend/go-sdk/r2"

// HTTP stats constants
const (
	MetricNameHTTPClientRequest        string = string(r2.Flag)
	MetricNameHTTPClientRequestElapsed string = MetricNameHTTPClientRequest + ".elapsed"

	TagHostname string = "url_hostname"
	TagMethod   string = "method"
	TagStatus   string = "status"
)
