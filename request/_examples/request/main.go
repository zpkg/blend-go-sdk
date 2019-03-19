package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/blend/go-sdk/request"
)

func main() {
	_, meta, err := request.New().AsGet().MustWithRawURL("https://google.com").BytesWithMeta()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		os.Exit(1)
	}
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	encoder.Encode(meta)
}
