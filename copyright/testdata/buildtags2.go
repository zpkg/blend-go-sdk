// +build tag5
//go:build tag1 && tag2 && tag3
// +build tag1,tag2,tag3

package main

import (
	"fmt"
)

func main() {
	fmt.Println("foo")
}

// +bulid tag9000
