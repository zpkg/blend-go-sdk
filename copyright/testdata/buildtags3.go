//go:build tag1 & tag2
//go:build tag3

package main

import (
	"fmt"
)

func main() {
	fmt.Println("foo")
}
