/*

Copyright (c) 2022 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package collections

// Error is an error string.
type Error string

// Error implements error.
func (e Error) Error() string { return string(e) }

// Labels is a loose type alias to map[string]string
type Labels = map[string]string

// Vars is a loose type alias to map[string]string
type Vars = map[string]interface{}

// Any is a loose type alias to interface{}.
type Any = interface{}
