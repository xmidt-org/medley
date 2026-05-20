// SPDX-FileCopyrightText: 2026 Comcast Cable Communications Management, LLC
// SPDX-License-Identifier: Apache-2.0

package medley

import (
	"hash"
	"io"

	"github.com/xmidt-org/medley/internal"
)

// HashResult defines the types of results of hashing supported
// by the stdlib.
type HashResult interface {
	uint32 | uint64
}

// Hash is medley's analog for hash.Hash. It behaves exactly like hash.Hash.
// The Sum32/Sum64 methods are replaced with the Value method.
//
// Hash also implements io.StringWriter, allowing strings to be written directly
// with no additional allocation, and io.ByteWriter. No errors are ever returned
// via these methods, just as with hash.Hash.
type Hash[HR HashResult] interface {
	hash.Hash
	io.StringWriter
	io.ByteWriter

	// Value is the method that returns Sum32 or Sum64, depending on the result type.
	// This method normalizes the hash.Hash interface, removing the need to have separate
	// code for hash.Hash32 and hash.Hash64.
	Value() HR
}

// hash32 is the internal adapter that implements Hash[uint32]
type hash32 struct {
	hash.Hash32
}

func (h32 *hash32) Value() uint32 {
	return h32.Sum32()
}

func (h32 *hash32) WriteString(v string) (int, error) {
	return h32.Write(
		internal.UnsafeBytes(v),
	)
}

func (h32 *hash32) WriteByte(c byte) error {
	var buf [1]byte
	buf[0] = c
	_, err := h32.Write(buf[:])
	return err
}

// AsHash32 converts a hash.Hash32 into a medley 32-bit Hash object.
func AsHash32(h32 hash.Hash32) Hash[uint32] {
	return &hash32{
		Hash32: h32,
	}
}

// AsConstructor32 converts a 32-bit hash constructor into a medley constructor.
// If ctor32 is nil, this function immediately panics.
func AsConstructor32(ctor32 func() hash.Hash32) func() Hash[uint32] {
	if ctor32 == nil {
		panic("a 32-bit hash constructor is required")
	}

	return func() Hash[uint32] {
		return AsHash32(ctor32())
	}
}

// hash64 is the internal adapter that implements Hash[uint64]
type hash64 struct {
	hash.Hash64
}

func (h64 *hash64) Value() uint64 {
	return h64.Sum64()
}

func (h64 *hash64) WriteString(v string) (int, error) {
	return h64.Write(
		internal.UnsafeBytes(v),
	)
}

func (h64 *hash64) WriteByte(c byte) error {
	var buf [1]byte
	buf[0] = c
	_, err := h64.Write(buf[:])
	return err
}

// AsHash64 converts a hash.Hash64 into a medley 64-bit Hash object.
func AsHash64(h64 hash.Hash64) Hash[uint64] {
	return &hash64{
		Hash64: h64,
	}
}

// AsConstructor64 converts a 64-bit hash constructor into a medley constructor.
// If ctor64 is nil, this function immediately panics.
func AsConstructor64(ctor64 func() hash.Hash64) func() Hash[uint64] {
	if ctor64 == nil {
		panic("a 64-bit hash constructor is required")
	}

	return func() Hash[uint64] {
		return AsHash64(ctor64())
	}
}
