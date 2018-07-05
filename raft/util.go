package raft

import (
	crand "crypto/rand"
	"fmt"
	"math"
	"math/big"
	"math/rand"
	"time"
)

var (
	_randSource = rand.NewSource(newSeed())
)

// Backoff is used to compute an exponential backoff
// duration. Base time is scaled by the current round,
// up to some maximum scale factor.
func Backoff(base time.Duration, backoffIndex int32) time.Duration {
	power := backoffIndex
	for power > 2 {
		base *= 2
		power--
	}
	return base
}

// RandomTimeout returns a value that is between the minVal and 3x minVal.
// i.e. it is minVal + ([0, 2 * minVal])
func RandomTimeout(minVal time.Duration) time.Duration {
	if minVal == 0 {
		return minVal
	}

	randomValue := time.Duration(_randSource.Int63())
	return minVal + (randomValue % (2 * minVal))
}

// min returns the minimum.
func min(a, b uint64) uint64 {
	if a <= b {
		return a
	}
	return b
}

// max returns the maximum.
func max(a, b uint64) uint64 {
	if a >= b {
		return a
	}
	return b
}

// returns an int64 from a crypto random source
// can be used to seed a source for a math/rand.
func newSeed() int64 {
	r, err := crand.Int(crand.Reader, big.NewInt(math.MaxInt64))
	if err != nil {
		panic(fmt.Errorf("failed to read random bytes: %v", err))
	}
	return r.Int64()
}
