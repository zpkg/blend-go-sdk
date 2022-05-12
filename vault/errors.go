/*

Copyright (c) 2022 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package vault

import "github.com/blend/go-sdk/ex"

// Common error codes.
const (
	ErrNotFound     ex.Class = "vault; not found"
	ErrUnauthorized ex.Class = "vault; not authorized"
	ErrServerError  ex.Class = "vault; remote error"
)
