// SPDX-FileCopyrightText: 2025 Comcast Cable Communications Management, LLC
// SPDX-License-Identifier: Apache-2.0

package consistent

import (
	"bytes"
	"strconv"

	"github.com/xmidt-org/medley"
)

// hasher implements all the low-level hashing logic for hash Rings.
type hasher[S medley.Service] struct {
	vnodes        int
	alg           medley.Algorithm
	serviceHasher medley.ServiceHasher[S]
}

// sum64 uses this hasher's algorithm to compute the hash token for
// the given object.
func (h hasher[S]) sum64(object []byte) uint64 {
	return h.alg.Sum64Bytes(object)
}

// ringSize returns the total number of nodes required to store the given
// number of services.
func (h hasher[S]) ringSize(serviceCount int) int {
	return h.vnodes * serviceCount
}

// base computes the hash bytes for a service used as the base
// for each computed token.
func (h hasher[S]) base(service S) []byte {
	var b bytes.Buffer
	h.serviceHasher(&b, service)
	return b.Bytes()
}

// serviceNodes computes the individual ring nodes for a single service.
func (h hasher[S]) serviceNodes(svc S) (snodes nodes[S]) {
	snodes = make(nodes[S], 0, h.vnodes)

	var (
		hash = h.alg.New64()
		base = h.base(svc)

		// a stack-allocated prefixBuffer to minimize allocations for the prefix bytes
		prefixBuffer [8]byte

		// prefix can grow beyond the initial buffer, but it's unlikely
		prefix = prefixBuffer[:]
	)

	for increment := int64(0); increment < int64(h.vnodes); increment++ {
		hash.Reset()
		prefix = strconv.AppendInt(prefix[:0], increment, 10)
		prefix = append(prefix, '=')
		hash.Write(prefix)
		hash.Write(base)

		snodes = append(snodes, &node[S]{token: hash.Sum64(), service: svc})
	}

	return
}
