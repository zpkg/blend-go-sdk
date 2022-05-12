/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package stringutil

import "strings"

// CSV produces a csv from a given set of values.
func CSV(values []string) string {
	return strings.Join(values, ",")
}
