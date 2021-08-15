/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package consistenthash

import (
	"fmt"
	"sort"
	"sync"
)

const (
	// DefaultReplicas is the default number of replicas.
	DefaultReplicas = 16
)

// New returns a new consistent hash.
func New(opts ...Option) *ConsistentHash {
	var ch ConsistentHash
	for _, opt := range opts {
		opt(&ch)
	}
	return &ch
}

// Option mutates a consistent hash.
type Option func(*ConsistentHash)

// OptBuckets sets the buckets list.
func OptBuckets(buckets ...string) Option {
	return func(ch *ConsistentHash) {
		for _, bucket := range buckets {
			ch.AddBucket(bucket)
		}
	}
}

// OptReplicas sets the virtual replica count.
//
// More virtual replicas can help with making item assignments
// more uniform, but the tradeoff is every operation takes a little
// longer as log2 of the number of buckets times the number of virtual replicas.
func OptReplicas(replicas int) Option {
	return func(ch *ConsistentHash) { ch.replicas = replicas }
}

// OptHashFunction sets the hash function.
//
// The default hash function is `consistenthash.StableHash` which uses
// a stable crc64 hash function to preserve ordering between process restarts.
func OptHashFunction(hashFunction HashFunction) Option {
	return func(ch *ConsistentHash) { ch.hashFunction = hashFunction }
}

// ConsistentHash creates hashed assignments for each bucket.
type ConsistentHash struct {
	mu sync.RWMutex

	replicas     int
	buckets      map[string]struct{}
	hashFunction HashFunction
	hashring     []hashedString
}

// ReplicasOrDefault is the default number of replicas.
func (ch *ConsistentHash) ReplicasOrDefault() int {
	if ch.replicas > 0 {
		return ch.replicas
	}
	return DefaultReplicas
}

// Buckets returns the buckets.
func (ch *ConsistentHash) Buckets() (buckets []string) {
	for bucket := range ch.buckets {
		buckets = append(buckets, bucket)
	}
	sort.Strings(buckets)
	return
}

// HashFunctionOrDefault returns the provided hash function or a default.
func (ch *ConsistentHash) HashFunctionOrDefault() HashFunction {
	if ch.hashFunction != nil {
		return ch.hashFunction
	}
	return StableHash
}

// AddBucket adds a bucket to the consistent hash.
//
// If the new bucket does not exist on the hash ring the
// assignments mappings will be updated for each bucket including
// the newly added bucket.
//
// If the new bucket already exists on the hash ring
// no further action is taken.
func (ch *ConsistentHash) AddBucket(newBucket string) {
	ch.mu.Lock()
	defer ch.mu.Unlock()

	if ch.buckets == nil {
		ch.buckets = make(map[string]struct{})
	}
	if _, ok := ch.buckets[newBucket]; ok {
		return
	}
	ch.buckets[newBucket] = struct{}{}
	ch.insert(newBucket)
}

// RemoveBucket removes a bucket from the consistent hash, and returns
// a boolean indicating if the provided bucket was found.
//
// If the bucket exists on the hash ring, the bucket and its replicas are removed
// and the item assignments are updated for the remaining buckets.
//
// If the bucket does not exist on the ring, no further action is taken.
func (ch *ConsistentHash) RemoveBucket(toRemove string) (ok bool) {
	ch.mu.Lock()
	defer ch.mu.Unlock()

	if ch.buckets == nil {
		return
	}
	if _, ok = ch.buckets[toRemove]; !ok {
		return
	}
	delete(ch.buckets, toRemove)
	for x := 0; x < ch.ReplicasOrDefault(); x++ {
		index := ch.search(ch.bucketHashKey(toRemove, x))
		ch.hashring = append(ch.hashring[:index], ch.hashring[index+1:]...)
	}
	return
}

// Assignment returns the bucket assignment for a given item.
func (ch *ConsistentHash) Assignment(item string) (bucket string) {
	ch.mu.RLock()
	defer ch.mu.RUnlock()
	index := ch.search(item)
	if index >= len(ch.hashring) {
		index = 0
	}
	bucket = ch.hashring[index].Value
	return
}

// IsAssigned returns if a given bucket is assigned a given item.
func (ch *ConsistentHash) IsAssigned(bucket, item string) (ok bool) {
	ok = bucket == ch.Assignment(item)
	return
}

// Assignments returns the assignments for a given list of items organized
// by the name of the bucket, and an array of the assigned items.
func (ch *ConsistentHash) Assignments(items ...string) map[string][]string {
	output := make(map[string][]string)
	for _, item := range items {
		bucket := ch.Assignment(item)
		output[bucket] = append(output[bucket], item)
	}
	return output
}

//
// helpers
//

// insert inserts a hashring bucket.
//
// insertion uses heap push to sort the items on insert.
func (ch *ConsistentHash) insert(bucket string) {
	for x := 0; x < ch.ReplicasOrDefault(); x++ {
		ch.hashring = insertionSort(ch.hashring, hashedString{
			Value:    bucket,
			Replica:  x,
			Hashcode: ch.hashcode(ch.bucketHashKey(bucket, x)),
		})
	}
}

// search does a binary search for the first hashring index whose
// node hashcode is >= the hashcode of a given item.
func (ch *ConsistentHash) search(item string) (index int) {
	return sort.Search(len(ch.hashring), ch.searchFn(ch.hashcode(item)))
}

// searchFn returns a closure searching for a given hashcode.
func (ch *ConsistentHash) searchFn(hashcode uint64) func(index int) bool {
	return func(index int) bool {
		return ch.hashring[index].Hashcode >= hashcode
	}
}

// bucketHashKey formats a hash key for a given bucket virtual replica.
func (ch *ConsistentHash) bucketHashKey(bucket string, index int) string {
	return bucket + "|" + fmt.Sprintf("%02d", index)
}

// hashcode creates a hashcode for a given string
func (ch *ConsistentHash) hashcode(item string) uint64 {
	return ch.HashFunctionOrDefault()([]byte(item))
}

// hashedString is a bucket in the hashring
// that holds the hashcode, the bucket name (as Value)
// and the virtual replica index.
type hashedString struct {
	Hashcode uint64
	Value    string
	Replica  int
}
