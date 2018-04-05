package main

import (
	"fmt"
	"os"

	request "github.com/blendlabs/go-request"
	util "github.com/blendlabs/go-util"
)

func main() {
	_, meta, err := request.New().AsGet().WithURLf("https://google.com").BytesWithMeta()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		os.Exit(1)
	}
	fmt.Fprintf(os.Stdout, util.JSON.SerializePretty(meta, "", "  "))
	//fmt.Fprintf(os.Stdout, "%s\n", string(res))
}
