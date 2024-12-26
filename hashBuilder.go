package medley

import (
	"io"
	"math"
	"unsafe"
)

// HashBuilder is a Fluent Builder for hash values. Instances of this
// type may be created via NewHashBuilder. If created directly, then
// Use must be called before attempting any writes to the builder.
//
// This type is convenient for building hashes of complex types, such as structs.
// The HashBasicServiceTo function in this package gives an example of this usage.
type HashBuilder struct {
	dst   io.Writer
	sum64 func() uint64
	reset func()
	err   error
}

// NewHashBuilder creates a HashBuilder that writes to this given destination.
// Errors from the various WriteXXX methods are accumulated and available via Err().
// When an error occurs on the underlying destination, all subsequent writes are ignored.
//
// The given dst io.Writer may be a hash.Hash64. If so, then this builder's Sum64
// method may be used to extract the current hash value.
func NewHashBuilder(dst io.Writer) *HashBuilder {
	return new(HashBuilder).Use(dst)
}

// Use replaces this builder's underlying writer and resets error state. This method
// may be used to reuse a single builder for different hashes.
//
// If a HashBuilder is created directly, without using NewHashBuilder, this method
// is required to initialize the builder or writing will cause a panic.
func (hb *HashBuilder) Use(dst io.Writer) *HashBuilder {
	hb.dst = dst
	hb.err = nil

	// summer holds the behavior of anything that can compute a 64-bit hash value.
	// hash.Hash64 satisfies this interface, for example.
	type summer interface {
		Sum64() uint64
	}

	if s, ok := hb.dst.(summer); ok {
		hb.sum64 = s.Sum64
	}

	// resetter defines the behavior of an io.Writer that can be reset.
	type resetter interface {
		Reset()
	}

	if r, ok := hb.dst.(resetter); ok {
		hb.reset = r.Reset
	}

	return hb
}

// Err returns the first error that occurred in any Fluent Chain using this builder.
// When any error occurs, all subsequent WriteXXX methods do nothing.
func (hb *HashBuilder) Err() error {
	return hb.err
}

// CanSum64 returns true if the currently wrapped io.Writer can compute the
// current hash value, i.e. if it supplies a Sum64() uint64 method.
func (hb *HashBuilder) CanSum64() bool {
	return hb.sum64 != nil
}

// Sum64 returns the current, computed hash value. If the currently wrapped io.Writer
// does not provide a Sum64() uint64 method, this method returns zero (0).
//
// To determine if this method can compute the current hash value, use CanSum64.
func (hb *HashBuilder) Sum64() (v uint64) {
	if hb.sum64 != nil {
		v = hb.sum64()
	}

	return
}

// CanReset tests if Reset will actually do anything, i.e. if the currently wrapped
// io.Writer supplies a Reset method.
func (hb *HashBuilder) CanReset() bool {
	return hb.reset != nil
}

// Reset resets the underlying io.Writer to its initial state. If the currently wrapped
// io.Writer does not supply a Reset() method, this method does nothing.
func (hb *HashBuilder) Reset() {
	if hb.reset != nil {
		hb.reset()
	}
}

// Write writes the given bytes.
func (hb *HashBuilder) Write(v []byte) *HashBuilder {
	if hb.err == nil {
		_, hb.err = hb.dst.Write(v)
	}

	return hb
}

// WriteString writes the given string in a manner that does not require
// additional allocations.
func (hb *HashBuilder) WriteString(v string) *HashBuilder {
	if hb.err == nil && len(v) > 0 {
		_, hb.err = hb.dst.Write(
			unsafe.Slice(unsafe.StringData(v), len(v)),
		)
	}

	return hb
}

// WriteUint8 writes the given 8-bit unsigned integer.
func (hb *HashBuilder) WriteUint8(v uint8) *HashBuilder {
	if hb.err == nil {
		var buf [1]byte
		buf[0] = byte(v)
		_, hb.err = hb.dst.Write(buf[:])
	}

	return hb
}

// WriteUint16 writes the given 16-bit integer in big endian form.
func (hb *HashBuilder) WriteUint16(v uint16) *HashBuilder {
	if hb.err == nil {
		var buf [2]byte
		buf[0] = byte(v >> 8)
		buf[1] = byte(v)
		_, hb.err = hb.dst.Write(buf[:])
	}

	return hb
}

// WriteUint32 writes the given 32-bit integer in big endian form.
func (hb *HashBuilder) WriteUint32(v uint32) *HashBuilder {
	if hb.err == nil {
		var buf [4]byte
		buf[0] = byte(v >> 24)
		buf[1] = byte(v >> 16)
		buf[2] = byte(v >> 8)
		buf[3] = byte(v)
		_, hb.err = hb.dst.Write(buf[:])
	}

	return hb
}

// WriteUint64 writes the given 64-bit integer in big endian form.
func (hb *HashBuilder) WriteUint64(v uint64) *HashBuilder {
	if hb.err == nil {
		var buf [8]byte
		buf[0] = byte(v >> 56)
		buf[1] = byte(v >> 48)
		buf[2] = byte(v >> 40)
		buf[3] = byte(v >> 32)
		buf[4] = byte(v >> 24)
		buf[5] = byte(v >> 16)
		buf[6] = byte(v >> 8)
		buf[7] = byte(v)
		_, hb.err = hb.dst.Write(buf[:])
	}

	return hb
}

// WriteFloat32 writes a 32-bit float value in big endian form.
func (hb *HashBuilder) WriteFloat32(v float32) *HashBuilder {
	if hb.err == nil {
		hb = hb.WriteUint32(math.Float32bits(v))
	}

	return hb
}

// WriteFloat64 writes a 64-bit float value in big endian form.
func (hb *HashBuilder) WriteFloat64(v float64) *HashBuilder {
	if hb.err == nil {
		hb = hb.WriteUint64(math.Float64bits(v))
	}

	return hb
}
