// SPDX-FileCopyrightText: 2026 Comcast Cable Communications Management, LLC
// SPDX-License-Identifier: Apache-2.0

package medley

import (
	"hash/fnv"

	"github.com/spaolacci/murmur3"
	"github.com/xmidt-org/medley/internal"
)

// asSum produces a SumXX() function using a constructor. Useful when no package-level
// SumXX() function is provided.
func asSum[HR HashResult](ctor func() Hash[HR]) func([]byte) HR {
	return func(b []byte) HR {
		h := ctor()
		h.Write(b)
		return h.Value()
	}
}

// Algorithm is a hashing algorithm used by medley. An Algorithm is immutable and safe
// for concurrent use.
//
// The zero value for this type is invalid. Use NewAlgorithm to create an Algorithm.
type Algorithm[HR HashResult] struct {
	ctor func() Hash[HR]
	sum  func([]byte) HR
}

// New constructs a Hash object that can be used exactly like a hash.Hash.
func (alg *Algorithm[HR]) New() Hash[HR] {
	return alg.ctor()
}

// Sum produces a hash of a sequence of bytes. Most algorithms provide a a sum function
// that avoids some overhead of using the Hash32.
func (alg *Algorithm[HR]) Sum(b []byte) HR {
	return alg.sum(b)
}

// SumString produces a sum for a string. The string's bytes are obtained without
// a reallocation.
func (alg *Algorithm[HR]) SumString(v string) HR {
	return alg.sum(
		internal.UnsafeBytes(v),
	)
}

// NewAlgorithm constructs a medley algorithm of a particular result size. Algorithms are immutable
// and safe for concurrent access.
//
// The ctor function is required, and if not supplied this function immediately panics.
// Use the AsConstructor32 and AsConstructor64 functions to convert constructors in
// other packages, e.g. crc32.NewIEEE.
//
// The sum function is optional. Most hash packages provide a function with this signature
// to allow hashing a sequence of bytes without the overhead of constructing a Hash. If this
// sum function is nil, the returned Algorithm uses a sum function built in terms of the ctor.
func NewAlgorithm[HR HashResult](ctor func() Hash[HR], sum func([]byte) HR) *Algorithm[HR] {
	if ctor == nil {
		panic("a constructor is required to create an Algorithm")
	}

	alg := &Algorithm[HR]{
		ctor: ctor,
		sum:  sum,
	}

	if alg.sum == nil {
		alg.sum = asSum(ctor)
	}

	return alg
}

var default32 = NewAlgorithm(
	AsConstructor32(murmur3.New32),

	// We can't use the murmur3.Sum32 function right now because of:
	// https://github.com/spaolacci/murmur3/issues/34
	nil,
)

var default64 = NewAlgorithm(
	AsConstructor64(murmur3.New64),

	// We can't use the murmur3.Sum32 function right now because of:
	// https://github.com/spaolacci/murmur3/issues/34
	nil,
)

// Default32 returns medley's default 32-bit hashing algorithm, which is 32-bit murmur3
// with the default seed.
func Default32() *Algorithm[uint32] {
	return default32
}

// Default64 returns medley's default 64-bit hashing algorithm, which is 64-bit murmur3
// with the default seed.
func Default64() *Algorithm[uint64] {
	return default64
}

var fnv32 = NewAlgorithm(AsConstructor32(fnv.New32), nil)
var fnv32a = NewAlgorithm(AsConstructor32(fnv.New32a), nil)

// FNV32 returns the medley Algorithm for 32-bit fnv.
func FNV32() *Algorithm[uint32] { return fnv32 }

// FNV32a returns the medley Algorithm for 32-bit fnv-a.
func FNV32a() *Algorithm[uint32] { return fnv32a }

var fnv64 = NewAlgorithm(AsConstructor64(fnv.New64), nil)
var fnv64a = NewAlgorithm(AsConstructor64(fnv.New64a), nil)

// FNV64 returns the medley Algorithm for 64-bit fnv.
func FNV64() *Algorithm[uint64] { return fnv64 }

// FNV64a returns the medley Algorithm for 64-bit fnv-a.
func FNV64a() *Algorithm[uint64] { return fnv64a }
