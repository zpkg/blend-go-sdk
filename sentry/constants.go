/*

Copyright (c) 2022 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package sentry

import "github.com/zpkg/blend-go-sdk/logger"

// Constants
const (
	Platform     = "go"
	SDK          = "sentry.go"
	ListenerName = "sentry"
)

var (
	// DefaultListenerFlags are the default log flags to notify Sentry for
	DefaultListenerFlags = []string{logger.Error, logger.Fatal}
)
