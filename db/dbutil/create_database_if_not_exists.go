/*

Copyright (c) 2022 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package dbutil

import (
	"context"

	"github.com/zpkg/blend-go-sdk/db"
	"github.com/zpkg/blend-go-sdk/env"
	"github.com/zpkg/blend-go-sdk/ex"
)

// CreateDatabaseIfNotExists creates a databse if it doesn't exist.
//
// It will check if a given `serviceEnv` is prodlike, and if the database doesn't exist, and the `serviceEnv`
// is prodlike, an `ErrDatabaseDoesntExist` will be returned.
//
// If a given `serviceEnv` is not prodlike, the database will be created with a management connection.
func CreateDatabaseIfNotExists(ctx context.Context, serviceEnv, database string, opts ...db.Option) error {
	exists, err := DatabaseExists(ctx, database, opts...)
	if err != nil {
		return err
	}
	if !exists {
		if env.IsProdlike(serviceEnv) {
			return ex.New(ErrDatabaseDoesntExist, ex.OptMessagef("database: %s", database))
		}
		if err = CreateDatabase(ctx, database, opts...); err != nil {
			return err
		}
	}
	return nil
}
