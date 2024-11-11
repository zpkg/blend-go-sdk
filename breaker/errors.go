/*

Copyright (c) 2022 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package breaker

import "github.com/zpkg/blend-go-sdk/ex"

var (
	// ErrTooManyRequests is returned when the CB state is half open and the requests count is over the cb maxRequests
	ErrTooManyRequests ex.Class = "too many requests"
	// ErrOpenState is returned when the CB state is open
	ErrOpenState ex.Class = "circuit breaker is open"
)

// ErrIsOpen returns if the error is an ErrOpenState.
func ErrIsOpen(err error) bool {
	return ex.Is(err, ErrOpenState)
}

// ErrIsTooManyRequests returns if the error is an ErrTooManyRequests.
func ErrIsTooManyRequests(err error) bool {
	return ex.Is(err, ErrTooManyRequests)
}
