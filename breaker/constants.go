/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package breaker

import "time"

// Constants
const (
	DefaultClosedExpiryInterval = 5 * time.Second
	DefaultOpenExpiryInterval   = 60 * time.Second
	DefaultHalfOpenMaxActions   = 1
	DefaultConsecutiveFailures  = 5
)
