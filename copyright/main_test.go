/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package copyright

import (
	"testing"

	"github.com/blend/go-sdk/logger"
	"github.com/blend/go-sdk/testutil"
)

func TestMain(m *testing.M) {
	testutil.MarkUpdateGoldenFlag()

	testutil.New(
		m,
		testutil.OptLog(logger.All()),
	).Run()
}
