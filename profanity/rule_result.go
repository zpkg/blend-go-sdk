/*

Copyright (c) 2022 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package profanity

// RuleResult is a result from a rule.
type RuleResult struct {
	OK      bool
	File    string
	Line    int
	Message string
	Err     error
}
