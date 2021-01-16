/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT
license that can be found in the LICENSE file.

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
