// SPDX-FileCopyrightText: 2026 Comcast Cable Communications Management, LLC
// SPDX-License-Identifier: Apache-2.0

package consistent

import (
	"github.com/xmidt-org/medley"
)

// Ring is a 64-bit hash ring used for consistent hashing.
//
// The NearestXX() methods are safe for concurrent use. If Clear() is used,
// this Ring must be externally synchronized.
type Ring[V any] struct {
	sum   medley.Sum[uint64]
	nodes hashNodes[V]
}

// Nearest hashes an object and returns the nearest value on the ring.
// If this Ring is empty, it returns the zero value for V.
func (r *Ring[V]) Nearest(object []byte) (value V) {
	if r.nodes.Len() > 0 {
		value = r.nodes.nearest(
			r.sum(object),
		)
	}

	return
}

// NearestString hashes an object and returns the nearest value on the ring.
// If this Ring is empty, it returns the zero value for V.
func (r *Ring[V]) NearestString(object string) (value V) {
	if r.nodes.Len() > 0 {
		value = r.nodes.nearest(
			medley.SumString(r.sum, object),
		)
	}

	return
}

// Clear wipes out this Ring, zeroing out each node and setting the internal
// nodes to nil. After this method is called, the NearestXX methods will return
// the zero value for V. This method is idempotent.
//
// Using this method before a Ring is retired can result in significantly less gc pressure.
func (r *Ring[V]) Clear() {
	r.nodes.clear()
	r.sum = nil
	r.nodes = nil
}
