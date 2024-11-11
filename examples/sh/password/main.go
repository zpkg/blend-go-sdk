/*

Copyright (c) 2022 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package main

import (
	"fmt"

	"github.com/zpkg/blend-go-sdk/sh"
)

func main() {
	value, err := sh.Password("first? ")
	if err != nil {
		sh.Fatal(err)
	}
	fmt.Println("entered", value)

	value, err = sh.Passwordf("%s? ", "second")
	if err != nil {
		sh.Fatal(err)
	}
	fmt.Println("entered", value)
}
