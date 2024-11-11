/*

Copyright (c) 2022 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package collections

import (
	"fmt"
	"testing"

	"github.com/zpkg/blend-go-sdk/assert"
)

func Test_BatchIterator(t *testing.T) {
	its := assert.New(t)

	bi := &BatchIterator{BatchSize: 100}
	its.False(bi.HasNext())
	its.Empty(bi.Next())

	bi = &BatchIterator{Items: generateBatchItems(10)}
	its.True(bi.HasNext())
	its.Empty(bi.Next())

	// handle edge case where somehow the cursor gets set beyond the
	// last element of the items.
	bi = &BatchIterator{Items: generateBatchItems(10), Cursor: 15}
	its.False(bi.HasNext())
	its.Empty(bi.Next())

	bi = &BatchIterator{Items: generateBatchItems(10), BatchSize: 100}
	its.True(bi.HasNext())
	its.Len(bi.Next(), 10)
	its.False(bi.HasNext())

	bi = &BatchIterator{Items: generateBatchItems(100), BatchSize: 10}
	for x := 0; x < 10; x++ {
		its.True(bi.HasNext())
		its.Len(bi.Next(), 10, fmt.Sprintf("failed on pass %d", x))
	}
	its.False(bi.HasNext())

	bi = &BatchIterator{Items: generateBatchItems(105), BatchSize: 10}
	for x := 0; x < 10; x++ {
		its.True(bi.HasNext())
		its.Len(bi.Next(), 10, fmt.Sprintf("failed on pass %d", x))
	}
	its.True(bi.HasNext())
	its.Len(bi.Next(), 5)
	its.False(bi.HasNext())
}

func generateBatchItems(count int) []interface{} {
	var output []interface{}
	for x := 0; x < count; x++ {
		output = append(output, fmt.Sprint(x))
	}
	return output
}
