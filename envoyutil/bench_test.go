/*

Copyright (c) 2022 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package envoyutil_test

import (
	"testing"

	"github.com/zpkg/blend-go-sdk/envoyutil"
)

// BenchmarkParseXFCC tries to help determine a baseline of the speed for
// `envoyutil.ParseXFCC()` on "small" well-formed input.
//
// Can be run via: `go test -bench=.`
func BenchmarkParseXFCC(b *testing.B) {
	xfcc := "By=spiffe://cluster.local/ns/blent/sa/echo;Hash=468ed33be74eee6556d90c0149c1309e9ba61d6425303443c0748a02dd8de688;Subject=10;URI=spiffe://cluster.local/ns/blent/sa/beep"
	for n := 0; n < b.N; n++ {
		_, _ = envoyutil.ParseXFCC(xfcc)
	}
}
