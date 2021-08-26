/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package ex

import (
	"errors"
	"fmt"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestIs(t *testing.T) {
	it := assert.New(t)

	stdLibErr := errors.New("sentinel")

	testCases := []struct {
		Err      interface{}
		Cause    error
		Expected bool
	}{
		{Err: Class("test class"), Cause: Class("test class"), Expected: true},
		{Err: New("test class"), Cause: Class("test class"), Expected: true},
		{Err: New("test class"), Cause: New("test class"), Expected: true},
		{Err: New(stdLibErr), Cause: stdLibErr, Expected: true},
		{Err: New(fmt.Errorf("outer err: %w", stdLibErr)), Cause: stdLibErr, Expected: true},
		{Err: Multi([]error{New("test class"), Class("not test class")}), Cause: Class("not test class"), Expected: true},
		{Err: Class("not test class"), Cause: New("test class"), Expected: false},
		{Err: New("test class"), Cause: New("not test class"), Expected: false},
		{Err: New("test class"), Cause: nil, Expected: false},
		{Err: fmt.Errorf("outer err: %w", stdLibErr), Cause: stdLibErr, Expected: true},
		{Err: nil, Cause: nil, Expected: false},
		{Err: nil, Cause: Class("test class"), Expected: false},
	}

	for index, tc := range testCases {
		it.Equal(tc.Expected, Is(tc.Err, tc.Cause), fmt.Sprintf("test case %d", index))
	}
}
