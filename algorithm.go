// SPDX-FileCopyrightText: 2025 Comcast Cable Communications Management, LLC
// SPDX-License-Identifier: Apache-2.0

package medley

import (
	"hash"
	"unsafe"

	"github.com/spaolacci/murmur3"
)

// Algorithm represents a hash algorithm which medley can use to implement
// service location.
type Algorithm struct {
	// New64 is the constructor for a Hash64 appropriate for this algorithm.
	// This field is required. If this field is unset, methods on this
	// Algorithm may panic.
	New64 func() hash.Hash64

	// Sum64 is this algorithm's simple function to compute a hash over a
	// byte slice. For many algorithms, this function will be more efficient
	// in simple cases.
	//
	// This field is not required. If not supplied, New64 will be used to
	// create a hash of the given bytes.
	Sum64 func([]byte) uint64
}

// Sumb64Bytes uses Sum64 to compute the hash of the given byte slice. If
// the Sum64 field isn't set, New64 is used to create a Hash64 and write
// the given bytes.
func (alg Algorithm) Sum64Bytes(v []byte) uint64 {
	if alg.Sum64 != nil {
		return alg.Sum64(v)
	}

	h := alg.New64()
	h.Write(v)
	return h.Sum64()
}

// Sum64String creates the hash of a string in a way that doesn't create
// unnecessary allocations.
//
// If Sum64 is set, that function is used to compute the hash. Otherwise,
// New64 is used to create a Hash64 and write the string's bytes.
func (alg Algorithm) Sum64String(v string) uint64 {
	return alg.Sum64Bytes(
		unsafe.Slice(unsafe.StringData(v), len(v)),
	)
}

// DefaultAlgorithm returns the default hash algorithm for medley.
// The returned object uses the murmur3 algorithm. The specific
// implementation is github.com/spaolacci/murmur3.
func DefaultAlgorithm() Algorithm {
	return Algorithm{
		New64: murmur3.New64,
		Sum64: murmur3.Sum64,
	}
}
