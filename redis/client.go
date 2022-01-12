/*

Copyright (c) 2022 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package redis

import (
	"context"
	"io"
)

// Client is the basic interface that a redis client should implement.
type Client interface {
	io.Closer
	Do(ctx context.Context, out interface{}, command string, args ...string) error
}
