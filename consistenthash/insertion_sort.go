/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package consistenthash

import "sort"

// insertionSort inserts an hashed string item by binary searching for the index
// which would satisfy the overall "sorted" status of the ring.
func insertionSort(ring []hashedString, item hashedString) []hashedString {
	destination := sort.Search(len(ring), func(index int) bool {
		return ring[index].Hashcode >= item.Hashcode
	})
	ring = append(ring, hashedString{})
	copy(ring[destination+1:], ring[destination:])
	ring[destination] = item
	return ring
}
