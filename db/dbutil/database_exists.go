/*

Copyright (c) 2022 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package dbutil

import (
	"context"

	"github.com/zpkg/blend-go-sdk/db"
)

// DatabaseExists returns if a database exists or not.
func DatabaseExists(ctx context.Context, name string, opts ...db.Option) (exists bool, err error) {
	var conn *db.Connection
	defer func() {
		err = db.PoolCloseFinalizer(conn, err)
	}()

	conn, err = OpenManagementConnection(opts...)
	if err != nil {
		return
	}

	exists, err = conn.QueryContext(ctx, "SELECT 1 FROM pg_database WHERE datname = $1", name).Any()
	return
}
