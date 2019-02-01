package main

import (
	"fmt"
	"os"
	"time"

	"github.com/blend/go-sdk/r2"
)

func main() {
	_, err := r2.New("https://google.com",
		r2.WithMethodGet(),
		r2.WithTimeout(500*time.Millisecond),
		r2.WithHeader("X-Sent-By", "go-sdk/request2"),
		r2.WithCookieValue("ssid", "baileydog01"),
	).CopyTo(os.Stdout)

	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
