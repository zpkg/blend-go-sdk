package main

import (
	"fmt"
	"os"
	"time"

	"github.com/blend/go-sdk/r2"
)

func main() {
	_, err := r2.New("https://google.com/robots.txt",
		r2.OptGet(),
		r2.OptTimeout(500*time.Millisecond),
		r2.OptHeaderValue("X-Sent-By", "go-sdk/request2"),
		r2.OptCookieValue("r2-ray-id", "baileydog01"),
	).CopyTo(os.Stdout)

	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
