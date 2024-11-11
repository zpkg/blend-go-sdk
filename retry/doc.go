/*

Copyright (c) 2022 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*Package retry implements a generic retry provider.

Basic Example

You can use the retry provider to wrap a potentially flaky call:

	res, err := http.Get("https://google.com")
	...
	res, err := retry.Retry(ctx, func(_ context.Context, args interface{}) (interface{}, error) {
		return http.Get(args.(string))
	}, "https://google.com") // note: the return type here is `(interface{}, error)`

You can also add additional parameters to the retry:

	res, err := retry.Retry(ctx, func(_ context.Context, args interface{}) (interface{}, error) {
		return http.Get(args.(string))
	}, "https://google.com", retry.OptMaxAttempts(10), retry.OptExponentialBackoff(time.Second))

*/
package retry // import "github.com/zpkg/blend-go-sdk/retry"
