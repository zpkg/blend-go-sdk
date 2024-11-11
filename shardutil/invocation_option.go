/*

Copyright (c) 2022 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package shardutil

import (
	"database/sql"
	"fmt"

	"github.com/zpkg/blend-go-sdk/db"
)

// InvocationOption is an option that returns an invocation option based on a partition index.
type InvocationOption func(partitionIndex int) db.InvocationOption

// OptTxs returns a base manager invocation option that
// parameterizes the transaction per invocation based on an array of transactions.
func OptTxs(txns ...*sql.Tx) InvocationOption {
	return func(partitionIndex int) db.InvocationOption {
		return func(i *db.Invocation) {
			db.OptTx(txns[partitionIndex])(i)
		}
	}
}

// OptLabel sets a label for invocations.
func OptLabel(label string) InvocationOption {
	return func(_ int) db.InvocationOption { return db.OptLabel(label) }
}

// OptPartitionLabel sets a label for invocations.
func OptPartitionLabel(label string) InvocationOption {
	return func(partitionIndex int) db.InvocationOption {
		return db.OptLabel(fmt.Sprintf("%s_%d", label, partitionIndex))
	}
}
