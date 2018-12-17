package main

import (
	"fmt"

	"github.com/blend/go-sdk/sh"
)

func main() {
	value, err := sh.Prompt("first? ")
	if err != nil {
		sh.Fatal(err)
	}
	fmt.Println("entered", value)
}
