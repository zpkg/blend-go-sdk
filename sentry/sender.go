/*

Copyright (c) 2022 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package sentry

import (
	"context"

	"github.com/zpkg/blend-go-sdk/logger"
)

// Sender is the type that the sentry client ascribes to.
type Sender interface {
	Notify(context.Context, logger.ErrorEvent)
}
