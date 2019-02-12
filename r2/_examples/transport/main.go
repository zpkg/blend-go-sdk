package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/blend/go-sdk/r2"
)

func main() {
	//create external transport reference
	transport := &http.Transport{}

	// pass to the request
	req, err := r2.New("https://google.com/robots.txt", r2.OptTransport(transport))
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(1)
	}

	var res *http.Response
	for x := 0; x < 10; x++ {
		res, err = r2.Close(req, nil)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
		} else {
			fmt.Fprintf(os.Stdout, "%v %v %v\n", res.StatusCode, res.Status, res.ContentLength)
		}
		time.Sleep(500 * time.Millisecond)
	}
	fmt.Println("Done")
}
