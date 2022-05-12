/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package configutil

// StringPtr returns a StringSource for a given string pointer.
func StringPtr(value *string) StringSource {
	if value == nil || *value == "" {
		return String("")
	}
	return String(*value)
}
