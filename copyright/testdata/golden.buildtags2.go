// +build tag5
//go:build tag1 && tag2 && tag3
// +build tag1,tag2,tag3

/*

Copyright (c) 2001 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package main

import (
	"fmt"
)

func main() {
	fmt.Println("foo")
}

// +bulid tag9000
