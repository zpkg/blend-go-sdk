r2
==

This is meant to be an experiment for an options based api for sending request. It is not stable and should only be used on an experimental basis.


## Usage

R2 uses a different paradigm from `go-sdk/request`; instead of chaining calls with a "fluent" api, options can be provided in a variadic list. This lets users extend the possible options as necessary.

## Example

```golang
package main

import (
	"fmt"
	"os"
	"time"

	"github.com/blend/go-sdk/r2"
)

func CustomOption() r2.Option {
	return func(r *r2.Request) {
		r.Client.Timeout = 10 * time.Millisecond
	}
}

func main() {
	_, err := r2.New("https://google.com",
		r2.Get(),
		r2.Timeout(500*time.Millisecond),
		r2.Header("X-Sent-By", "go-sdk/request2"),
		r2.CookieValue("ssid", "baileydog01"),
		CustomOption(),
	).CopyTo(os.Stdout)

	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
```