/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package db

// StatementInterceptor is an interceptor for statements.
type StatementInterceptor func(statementID, statement string) string
