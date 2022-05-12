/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

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
