/*

Copyright (c) 2022 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package main

import (
	"log"

	"github.com/zpkg/blend-go-sdk/ex"
)

func displayError(err error) error {
	asMulti, ok := err.(*multiError)
	if !ok {
		return err
	}

	if asMulti == nil || len(asMulti.Errors) == 0 {
		return ex.New("Expected a non-nil / non-empty error")
	}

	log.Println("Error(s):")
	for _, err := range asMulti.Errors {
		log.Printf("- %#v\n", err)
	}

	return nil
}
