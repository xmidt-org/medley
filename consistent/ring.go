// SPDX-FileCopyrightText: 2026 Comcast Cable Communications Management, LLC
// SPDX-License-Identifier: Apache-2.0

package consistent

import (
	"github.com/xmidt-org/medley"
)

// Ring is a 64-bit hash ring used for consistent hashing.
type Ring[V any] struct {
	alg   *medley.Algorithm[uint64]
	nodes hashNodes[V]
}

// Nearest hashes an object and returns the nearest value on the ring.
// If this Ring is empty, it returns the zero value for V.
//
// Nearest by itself is safe for concurrent access and does not mutate this Ring.
// However, Clear mutates this Ring and thus calling code must provide synchronization
// if Clear is used.
func (r *Ring[V]) Nearest(object medley.Object) (value V) {
	if r.nodes.Len() > 0 {
		value = r.nodes.nearest(
			r.alg.SumObject(object),
		)
	}

	return
}

// Clear wipes out this Ring, zeroing out each node and setting the internal
// nodes to nil. After this method is called, Nearest will return the zero value for V.
// This method is idempotent.
//
// Using this method before a Ring is retired can result in significantly less gc pressure.
//
// This method mutates this Ring, and if used must be synchronized with other calls
// to the same Ring.
func (r *Ring[V]) Clear() {
	r.nodes.clear()
	r.alg = nil
	r.nodes = nil
}
