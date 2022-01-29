/*

Copyright (c) 2022 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package collections

type listNode struct {
	Next     *listNode
	Previous *listNode
	Value    interface{}
}

// NewLinkedList returns a new linked list instance.
func NewLinkedList() *LinkedList {
	return &LinkedList{}
}

// NewLinkedListFromValues creates a linked list out of a slice.
func NewLinkedListFromValues(values []interface{}) *LinkedList {
	list := new(LinkedList)
	for _, v := range values {
		list.Enqueue(v)
	}
	return list
}

// LinkedList is an implementation of a fifo buffer using nodes and poitners.
// Remarks; it is not threadsafe. It is constant(ish) time in all ops.
type LinkedList struct {
	head   *listNode
	tail   *listNode
	length int
}

// Len returns the length of the queue in constant time.
func (q *LinkedList) Len() int {
	return q.length
}

// Enqueue adds a new value to the queue.
func (q *LinkedList) Enqueue(value interface{}) {
	node := &listNode{Value: value}

	if q.head == nil { //the queue is empty, that is to say head is nil
		q.head = node
		q.tail = node
	} else { //the queue is not empty, we have a (valid) tail pointer
		q.tail.Previous = node
		node.Next = q.tail
		q.tail = node
	}

	q.length++
}

// Dequeue removes an item from the front of the queue and returns it.
func (q *LinkedList) Dequeue() interface{} {
	if q.head == nil {
		return nil
	}

	headValue := q.head.Value

	if q.length == 1 && q.head == q.tail {
		q.head = nil
		q.tail = nil
	} else {
		q.head = q.head.Previous
		if q.head != nil {
			q.head.Next = nil
		}
	}

	q.length--
	return headValue
}

// DequeueBack pops the _last_ element off the linked list.
func (q *LinkedList) DequeueBack() interface{} {
	if q.tail == nil {
		return nil
	}
	tailValue := q.tail.Value

	if q.length == 1 {
		q.head = nil
		q.tail = nil
	} else {
		q.tail = q.tail.Next
		if q.tail != nil {
			q.tail.Previous = nil
		}
	}

	q.length--
	return tailValue
}

// Peek returns the first element of the queue but does not remove it.
func (q *LinkedList) Peek() interface{} {
	if q.head == nil {
		return nil
	}
	return q.head.Value
}

// PeekBack returns the last element of the queue.
func (q *LinkedList) PeekBack() interface{} {
	if q.tail == nil {
		return nil
	}
	return q.tail.Value
}

// Clear clears the linked list.
func (q *LinkedList) Clear() {
	q.tail = nil
	q.head = nil
	q.length = 0
}

// Drain calls the consumer for each element of the linked list.
func (q *LinkedList) Drain() []interface{} {
	if q.head == nil {
		return nil
	}

	contents := make([]interface{}, q.length)
	nodePtr := q.head
	var index int
	for nodePtr != nil {
		contents[index] = nodePtr.Value
		nodePtr = nodePtr.Previous
		index++
	}
	q.tail = nil
	q.head = nil
	q.length = 0
	return contents
}

// Each calls the consumer for each element of the linked list.
func (q *LinkedList) Each(consumer func(value interface{})) {
	if q.head == nil {
		return
	}

	nodePtr := q.head
	for nodePtr != nil {
		consumer(nodePtr.Value)
		nodePtr = nodePtr.Previous
	}
}

// Consume calls the consumer for each element of the linked list, removing it.
func (q *LinkedList) Consume(consumer func(value interface{})) {
	if q.head == nil {
		return
	}

	nodePtr := q.head
	for nodePtr != nil {
		consumer(nodePtr.Value)
		nodePtr = nodePtr.Previous
	}
	q.tail = nil
	q.head = nil
	q.length = 0
}

// EachUntil calls the consumer for each element of the linked list, but can abort.
func (q *LinkedList) EachUntil(consumer func(value interface{}) bool) {
	if q.head == nil {
		return
	}

	nodePtr := q.head
	for nodePtr != nil {
		if !consumer(nodePtr.Value) {
			return
		}
		nodePtr = nodePtr.Previous
	}
}

// ReverseEachUntil calls the consumer for each element of the linked list, but can abort.
func (q *LinkedList) ReverseEachUntil(consumer func(value interface{}) bool) {
	if q.head == nil {
		return
	}

	nodePtr := q.tail
	for nodePtr != nil {
		if !consumer(nodePtr.Value) {
			return
		}
		nodePtr = nodePtr.Next
	}
}

// Contents returns the full contents of the queue as a slice.
func (q *LinkedList) Contents() []interface{} {
	if q.head == nil {
		return []interface{}{}
	}

	values := []interface{}{}
	nodePtr := q.head
	for nodePtr != nil {
		values = append(values, nodePtr.Value)
		nodePtr = nodePtr.Previous
	}
	return values
}
