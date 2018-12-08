package main

import (
	"fmt"
	"os"

	"github.com/blend/go-sdk/sh"
)

func main() {
	if err := sh.Pipe(sh.C("yes", "head -n 5")...); err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}
}
