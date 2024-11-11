/*

Copyright (c) 2022 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*
Package migration provides helpers for writing rerunnable database migrations.

These are built around Suites, which are sets of Groups that execute within a transaction, those Groups are composed of Steps, which are a Guard and an Action.
*/
package migration // import "github.com/zpkg/blend-go-sdk/db/migration"
