/*

Copyright (c) 2022 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package statsd

import "github.com/blend/go-sdk/stats"

// Tag is an alias / wrapper to stats.Tag
func Tag(k, v string) string {
	return stats.Tag(k, v)
}
