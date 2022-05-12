/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package redis

import "context"

// Assert `MockClient` implements client.
var (
	_ Client = (*MockClient)(nil)
)

// MockClient is a mocked client.
type MockClient struct {
	Ops chan MockClientOp
}

// Do applies a command.
func (mc *MockClient) Do(_ context.Context, out interface{}, op string, args ...string) error {
	mc.Ops <- MockClientOp{Out: out, Op: op, Args: args}
	return nil
}

// Close closes the mock client.
func (mc *MockClient) Close() error { return nil }

// MockClientOp is a mocked client op.
type MockClientOp struct {
	Out  interface{}
	Op   string
	Args []string
}
