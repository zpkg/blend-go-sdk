/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*
Package consistenthash contains helpers for mapping items to buckets using consistent hashing.

The underlying hash function is `crc32.ChecksumIEEE` but that can be customized.

Consistent hash (the result of `New(...)`) is safe to use from multiple goroutines and
will use a mutex to synchronize changes to internal state.

Example usage:

    ch := consistenthash.New(
		consistenthash.OptBuckets("worker-0", "worker-1", "worker-2"),
		consistenthash.OptItems(items...),
	)
	worker0Items := ch.Assignments("worker-0")

`worker0Items` will now hold just the items that were mapped to "worker-0".
*/
package consistenthash // import "github.com/blend/go-sdk/consistenthash"
