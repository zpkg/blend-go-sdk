package main

import (
	"fmt"
	"os"

	"github.com/blend/go-sdk/request"
	"github.com/blend/go-sdk/util"
)

func main() {
	_, meta, err := request.New().AsGet().WithRawURLf("https://google.com").BytesWithMeta()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		os.Exit(1)
	}
	fmt.Fprintf(os.Stdout, util.JSON.SerializePretty(meta, "", "  "))
}
