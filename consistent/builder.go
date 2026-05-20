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
//
// By default, a Builder uses DefaultVNodes number of vnodes and the murmur3
// hashing algorithm.
type Builder[O medley.Object, V any] struct {
	vnodes int
	ctor   medley.Constructor[uint64]
	sum    medley.Sum[uint64]
}

// VNodes sets the number of virtual nodes per value. If v is nonpositive,
// DefaultVNodes is used.
func (b *Builder[O, V]) VNodes(v int) *Builder[O, V] {
	b.vnodes = v
	return b
}

// Algorithm defines the hashing algorithm this Builder will use. The supplied constructor
// and sum function must agree, i.e. must come from the same package or implement the same
// algorithm.
//
// If ctor is nil, regardless of sum, this builder is set to the default algorithm, which
// is to use murmur3.
//
// If sum is nil but ctor is provided, a sum function using the constructor will be synthesized.
func (b *Builder[O, V]) Algorithm(ctor medley.Constructor[uint64], sum medley.Sum[uint64]) *Builder[O, V] {
	switch {
	case ctor == nil:
		// reset to default
		b.ctor = nil
		b.sum = nil

	case sum == nil:
		b.ctor = ctor
		b.sum = medley.AsSum(ctor)

	default:
		b.ctor = ctor
		b.sum = sum
	}

	return b
}

// algorithm returns the hash constructor and sum functions to use. If these have not been
// set, a default using murmur3 is used instead.
func (b *Builder[O, V]) algorithm() (medley.Constructor[uint64], medley.Sum[uint64]) {
	if b.ctor != nil && b.sum != nil {
		return b.ctor, b.sum
	}

	return medley.Default64()
}

// allocateRing creates an empty ring and computes the number of vnodes to use
// when building the ring. If n > 0, the ring's nodes will be
// preallocated.
//
// The returned medley Hash can be used to build the ring's nodes. The returned Ring
// will be configured with the algorithm's sum function.
func (b *Builder[O, V]) allocateRing(n int) (ring *Ring[V], hash medley.Hash[uint64], vnodes int) {
	ctor, sum := b.algorithm()
	hash = ctor()

	ring = &Ring[V]{
		sum: sum,
	}

	if vnodes = b.vnodes; vnodes <= 0 {
		vnodes = DefaultVNodes
	}

	if n > 0 {
		ring.nodes = allocateHashNodes[V](vnodes * n)
	}

	return
}

// Build constructs a hash ring over the supplied values using this builder's configuration.
// The sequence of values is taken as is. No deduplication is done.
//
// The N parameter is the number of expected tuples that the values sequence will return. N is used
// as a hint for preallocation. If N is positive, the ring will be preallocated with space for N values
// and all the required vnodes. If N is nonpositive, no preallocation is done.
//
// The state of this builder is retained after this method returns.
func (b *Builder[O, V]) Build(n int, values iter.Seq2[O, V]) *Ring[V] {
	var (
		ring, hash, vnodes = b.allocateRing(n)

		// create an initial buffer that is likely big enough in most cases
		tokenBuffer = make([]byte, 256)
	)

	for object, value := range values {
		// if ExpectedValues was set appropriately, growing won't do anything as we'll have enough space.
		// however, this cuts down the number of allocations in the case where no hint was given.
		ring.nodes = ring.nodes.grow(vnodes)

		for i := range uint64(vnodes) {
			// this identifier is equivalent to:
			// https://github.com/billhathaway/consistentHash/blob/master/consistentHash.go#L60
			hash.Reset()
			tokenBuffer = strconv.AppendUint(tokenBuffer[:0], i, 10) // monotonic integer
			tokenBuffer = append(tokenBuffer, '=')                   // github.com/billhathaway/consistentHash uses this delimiter
			tokenBuffer = append(tokenBuffer, object...)

			// based on benchmarking, one big write to the hash is faster
			// than individual writes.
			hash.Write(tokenBuffer)

			ring.nodes = ring.nodes.append(hash.Value(), value)
		}
	}

	ring.nodes.sort()
	return ring
}
