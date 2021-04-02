/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package pagerduty

// Reference is a generic api object reference type.
type Reference struct {
	ID      string        `json:"id"`
	Type    ReferenceType `json:"type"`
	Summary string        `json:"summary,omitempty"`
	Self    string        `json:"self,omitempty"`
	HTMLUrl string        `json:"html_url,omitempty"`
}
