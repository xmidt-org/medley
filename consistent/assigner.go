package consistent

import (
	"bytes"
	"hash"
	"strconv"

	"github.com/xmidt-org/medley"
)

// separator is the byte sequence written between a hash increment value and the node.
// This value is used for backward compatibility with https://github.com/billhathaway/consistentHash.
var separator = []byte("=")

// assigner generates a sequence of hash values for nodes.
//
// Note: Currently, this type generates a sequence of hash values
// that are backward compatible with https://github.com/billhathaway/consistentHash.
//
// See https://github.com/billhathaway/consistentHash/blob/master/consistentHash.go#L62
type assigner struct {
	hasher hash.Hash64

	node     bytes.Buffer
	index    int64
	indexBuf []byte
}

// newAssigner creates an assigner that uses a given algorithm.
// The returned assigner is set to an empty node.  To begin assigning
// hash values to nodes, use Reset followed by Next.
func newAssigner(alg medley.Algorithm) *assigner {
	return &assigner{
		hasher:   alg(),
		indexBuf: make([]byte, 6), // a starting point large enough to reduce allocations
	}
}

// reset initializes this assigner to begin a new sequence of hash values
// for a specific node
func (a *assigner) reset(n medley.Node) {
	a.node.Reset()
	a.index = 0

	// bytes.Buffer.Write always returns a nil error
	n.WriteTo(&a.node) //nolint:errcheck
}

// next generates the next hash value for the node passed to the last
// call to reset
func (a *assigner) next() uint64 {
	a.hasher.Reset()
	a.indexBuf = strconv.AppendInt(a.indexBuf[:0], a.index, 10)
	a.index++

	// hash.Hash64.Write is documented as never returning an error
	a.hasher.Write(a.indexBuf)     //nolint:errcheck
	a.hasher.Write(separator)      //nolint:errcheck
	a.hasher.Write(a.node.Bytes()) //nolint:errcheck
	return a.hasher.Sum64()
}
