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
type Builder[O medley.Object, V any] struct {
	vnodes int
	alg    *medley.Algorithm[uint64]
}

// VNodes sets the number of virtual nodes per value. If v is nonpositive,
// DefaultVNodes is used.
func (b *Builder[O, V]) VNodes(v int) *Builder[O, V] {
	b.vnodes = v
	return b
}

// Algorithm sets the Algorithm64 used to both generate the hash Ring and hash
// objects to lookup values on the Ring. By default, a Builder will use
// medley.DefaultAlgorithm().
func (b *Builder[O, V]) Algorithm(v *medley.Algorithm[uint64]) *Builder[O, V] {
	b.alg = v
	return b
}

// allocateRing creates an empty ring and computes the number of vnodes to use
// when building the ring. If n > 0, the ring's nodes will be
// preallocated.
func (b *Builder[O, V]) allocateRing(n int) (r *Ring[V], vnodes int) {
	r = new(Ring[V])

	if b.alg != nil {
		r.alg = b.alg
	} else {
		r.alg = medley.Default64()
	}

	if vnodes = b.vnodes; vnodes <= 0 {
		vnodes = DefaultVNodes
	}

	if n > 0 {
		r.nodes = allocateHashNodes[V](vnodes * n)
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
		r, vnodes = b.allocateRing(n)

		h = r.alg.New()

		// create an initial buffer that is likely big enough in most cases
		tokenBuffer = make([]byte, 256)
	)

	for object, value := range values {
		// if ExpectedValues was set appropriately, growing won't do anything as we'll have enough space.
		// however, this cuts down the number of allocations in the case where no hint was given.
		r.nodes = r.nodes.grow(vnodes)

		for i := range uint64(vnodes) {
			// this identifier is equivalent to:
			// https://github.com/billhathaway/consistentHash/blob/master/consistentHash.go#L60
			h.Reset()
			tokenBuffer = strconv.AppendUint(tokenBuffer[:0], i, 10) // monotonic integer
			tokenBuffer = append(tokenBuffer, '=')                   // github.com/billhathaway/consistentHash uses this delimiter
			tokenBuffer = medley.Append(tokenBuffer, object)

			// based on benchmarking, one big write to the hash is faster
			// than individual writes.
			h.Write(tokenBuffer)

			r.nodes = r.nodes.append(h.Value(), value)
		}
	}

	r.nodes.sort()
	return r
}
