package main

import (
	"fmt"
	"os"
	"time"

	"github.com/blend/go-sdk/request2"
)

func main() {
	_, err := request2.MustNew("https://google.com",
		request2.WithMethodGet(),
		request2.WithTimeout(500*time.Millisecond),
		request2.WithHeader("X-Sent-By", "go-sdk/request2"),
		request2.WithCookieValue("ssid", "baileydog01"),
	).CopyTo(os.Stdout)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
