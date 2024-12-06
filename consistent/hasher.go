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

// totalCount produces the total number of nodes needed for a ring
// given a count of services.
func (s hasher[S]) totalCount(serviceCount int) int {
	return serviceCount * s.vnodes
}

// base computes the hash bytes for a service used as the base
// for each computed token.
func (s hasher[S]) base(service S) []byte {
	var b bytes.Buffer
	s.serviceHasher(&b, service)
	return b.Bytes()
}

// serviceNodes computes the individual ring nodes for a single service.
func (s hasher[S]) serviceNodes(svc S) (snodes nodes[S]) {
	snodes = make(nodes[S], 0, s.vnodes)

	var (
		h    = s.alg.New64()
		base = s.base(svc)

		// a stack-allocated prefixBuffer to minimize allocations for the prefix bytes
		prefixBuffer [8]byte

		// prefix can grow beyond the initial buffer, but it's unlikely
		prefix = prefixBuffer[:]
	)

	for increment := int64(0); increment < int64(s.vnodes); increment++ {
		h.Reset()
		prefix = strconv.AppendInt(prefix[:0], increment, 10)
		prefix = append(prefix, '=')
		h.Write(prefix)
		h.Write(base)

		snodes = append(snodes, &node[S]{token: h.Sum64(), service: svc})
	}

	return
}

// hashString uses this sequence's configuration to hash the given string.
func (s hasher[S]) hashString(object string) uint64 {
	return medley.HashString(s.alg, object)
}
