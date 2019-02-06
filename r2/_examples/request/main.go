package main

import (
	"fmt"
	"os"
	"time"

	"github.com/blend/go-sdk/r2"
)

func main() {
	_, err := r2.New("https://google.com",
		r2.Get(),
		r2.Timeout(500*time.Millisecond),
		r2.HeaderValue("X-Sent-By", "go-sdk/request2"),
		r2.CookieValue("ssid", "baileydog01"),
	).CopyTo(os.Stdout)

	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
