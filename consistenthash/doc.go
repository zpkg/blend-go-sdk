/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*
Package consistenthash contains helpers for mapping items to buckets using consistent hashing.

The underlying hash function is `crc64.ChecksumIEEE` but that can be customized. The hash function
is seeded with a constant polynomial so that assignments are stable between process starts.

Consistent hash (the result of `New(...)`) is safe to use from multiple goroutines and
will use a mutex to synchronize changes to internal state.

Methods typically run in `log2(N*M)` time where N is the number of buckets and M is the
number of virtual replicas in use (wich defaults to 16). This is strictly worse than a typical
map, but avoids space issues with tracking item assignments individually.

Example usage:

    ch := consistenthash.New(
		consistenthash.OptBuckets("worker-0", "worker-1", "worker-2"),
	)
	// figure out which bucket an item maps to
	worker := ch.Assignment("item-0") // will yield `worker-0` or `worker-1` etc.

You can tune the number of virtual replicas to reduce the constant time hit of most operations
at the expense of bucket to item mapping uniformity.

Example setting the replicas:
    ch := consistenthash.New(
		consistenthash.OptBuckets("worker-0", "worker-1", "worker-2"),
		consistenthash.OptReplicas(5),
	)
	// figure out which bucket an item maps to
	worker := ch.Assignment("item-0") // will yield `worker-0` or `worker-1` etc.

*/
package consistenthash // import "github.com/blend/go-sdk/consistenthash"
