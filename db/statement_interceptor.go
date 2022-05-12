/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package db

// StatementInterceptor is an interceptor for statements.
type StatementInterceptor func(statementID, statement string) string
