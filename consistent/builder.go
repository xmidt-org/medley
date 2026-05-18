// SPDX-FileCopyrightText: 2026 Comcast Cable Communications Management, LLC
// SPDX-License-Identifier: Apache-2.0

package consistent

import (
	"iter"
	"strconv"

	"github.com/xmidt-org/medley"
)

const (
	// DefaultVNodes is the number of vnodes to use when no value is specified.
	DefaultVNodes = 200
)

// Builder is a Fluent Builder for consistent hash Rings. The zero value of this
// type is ready to use.
type Builder[V any] struct {
	vnodes         int
	alg            *medley.Algorithm[uint64]
	expectedValues int
}

// VNodes sets the number of virtual nodes per value. If v is nonpositive,
// DefaultVNodes is used.
func (b *Builder[V]) VNodes(v int) *Builder[V] {
	b.vnodes = v
	return b
}

// Alg64 sets the Algorithm64 used to both generate the hash Ring and hash
// objects to lookup values on the Ring. By default, a Builder will use
// medley.DefaultAlgorithm().
func (b *Builder[V]) Alg64(v *medley.Algorithm[uint64]) *Builder[V] {
	b.alg = v
	return b
}

// ExpectedValues gives the builder a hint as to the number of values in the sequence
// passed to Build. For example, if hashing a ring of 20 servers, use ExpectedValues(20)
// to have Build preallocate the ring.
func (b *Builder[V]) ExpectedValues(v int) *Builder[V] {
	b.expectedValues = v
	return b
}

// allocateRing creates an empty ring and computes the number of vnodes to use
// when building the ring. If expectedValues > 0, the ring's nodes will be
// preallocated.
func (b *Builder[V]) allocateRing() (r *Ring[V], vnodes int) {
	r = new(Ring[V])

	if b.alg != nil {
		r.alg = b.alg
	} else {
		r.alg = medley.Default64()
	}

	if vnodes = b.vnodes; vnodes <= 0 {
		vnodes = DefaultVNodes
	}

	if b.expectedValues > 0 {
		r.nodes = allocateHashNodes[V](vnodes * b.expectedValues)
	}

	return
}

// Build constructs a hash ring over the supplied values using this builder's configuration.
// If successful, this method returns a non-nil Ring and a nil error. If any error occurs,
// a nil Ring is returned along with that error.
//
// If possible, use ExpectedValues and supply the number of elements that the values sequence
// will return. This allows Build to preallocate nodes in the hash ring, rather than allocating
// on the fly as the value sequence is used.
//
// The state of this builder is retained after this method returns.
func (b *Builder[V]) Build(values iter.Seq2[medley.Object, V]) *Ring[V] {
	var (
		r, vnodes = b.allocateRing()

		h = r.alg.New()

		// we know this buffer is large enough to hold any uint16 in base 10
		iBuffer = make([]byte, 0, 5)

		// github.com/billhathaway/consistentHash uses this delimiter
		delimiter = []byte{'='}
	)

	for id, value := range values {
		// if ExpectedValues was set appropriately, growing won't do anything as we'll have enough space.
		// however, this cuts down the number of allocations in the case where no hint was given.
		r.nodes = r.nodes.grow(vnodes)

		for i := range uint64(vnodes) {
			// this identifier is equivalent to:
			// https://github.com/billhathaway/consistentHash/blob/master/consistentHash.go#L60
			h.Reset()
			h.Write(strconv.AppendUint(iBuffer, i, 10))
			h.Write(delimiter)
			id.ToHash(h)
			r.nodes = r.nodes.append(h.Value(), value)
		}
	}

	r.nodes.sort()
	return r
}
