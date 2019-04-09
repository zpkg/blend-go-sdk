package main

import (
	"os"

	"github.com/blend/go-sdk/fileutil"
)

func main() {
	go fileutil.Watch("file.yml", func(f *os.File) error {
		defer f.Close()
		println("file changed")
		return nil
	})

	select {}
}
