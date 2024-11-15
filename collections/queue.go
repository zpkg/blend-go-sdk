/*

Copyright (c) 2022 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package collections

// Queue is an interface for implementations of a FIFO buffer.
//
// With the DequeueBack method, you can also use a Queue as a stack.
type Queue interface {
	Len() int
	Enqueue(value interface{})
	Dequeue() interface{}
	DequeueBack() interface{}
	Peek() interface{}
	PeekBack() interface{}
	Drain() []interface{}
	Contents() []interface{}
	Clear()

	Consume(consumer func(value interface{}))
	Each(consumer func(value interface{}))
	EachUntil(consumer func(value interface{}) bool)
	ReverseEachUntil(consumer func(value interface{}) bool)
}
