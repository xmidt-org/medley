package medley

import "io"

// Key defines the behavior of hash keys.  All hash keys can write
// their own hashable bytes to an output sink.
type Key interface {
	io.WriterTo
}

// ComputeHash is a convenience for computing a single key's hash
// value using an algorithm.  In general, the hash object returned
// from an Algorithm should be reset and reused when computing
// many hash values.  This function is provided as a utility
// for test code and a convenience for tools that can query a hash.
func ComputeHash(k Key, alg Algorithm) uint64 {
	h := alg()

	// hash.Hash64 never returns an error from Write
	k.WriteTo(h) //nolint:errcheck

	return h.Sum64()
}

// Bytes is a Key which is a slice of bytes
type Bytes []byte

// WriteTo writes this Bytes key's contents to the given writer
func (b Bytes) WriteTo(w io.Writer) (int64, error) {
	c, err := w.Write([]byte(b))
	return int64(c), err
}

// String is a Key which is a golang string
type String string

// WriteTo writes this string's contents to the given writer.
// io.WriteString is used to optimize string writing where possible.
func (s String) WriteTo(w io.Writer) (int64, error) {
	c, err := io.WriteString(w, string(s))
	return int64(c), err
}
