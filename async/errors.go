/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package async

import "github.com/blend/go-sdk/ex"

// Errors
var (
	ErrCannotStart ex.Class = "cannot start; already started"
	ErrCannotStop  ex.Class = "cannot stop; already stopped"
)
