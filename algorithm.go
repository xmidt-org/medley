package medley

import (
	"hash"
	"unsafe"

	"github.com/spaolacci/murmur3"
)

// Algorithm represents a hash algorithm. This interface defines the
// hash behavior required by medley.
type Algorithm interface {
	// New64 creates a new 64-bit hasher.
	New64() hash.Hash64

	// Sum64 hashes the given bytes into a 64-bit integer.
	Sum64([]byte) uint64
}

// Murmur3 provides the murmur3 hash algorithm. The specific implementation
// is github.com/spaolacci/murmur3. This is the default underlying hashing
// algorithm used for medley.
type Murmur3 struct{}

func (Murmur3) New64() hash.Hash64    { return murmur3.New64() }
func (Murmur3) Sum64(v []byte) uint64 { return murmur3.Sum64(v) }

// HashString produces the hash of a string without performing extra
// allocations.
func HashString(v string, alg Algorithm) uint64 {
	return alg.Sum64(
		unsafe.Slice(unsafe.StringData(v), len(v)),
	)
}
