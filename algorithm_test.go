// SPDX-FileCopyrightText: 2026 Comcast Cable Communications Management, LLC
// SPDX-License-Identifier: Apache-2.0

package medley

import (
	"hash"
	"hash/fnv"
	"testing"

	"github.com/spaolacci/murmur3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type AlgorithmTestSuite[HR HashResult] struct {
	suite.Suite

	testString string
	testBytes  []byte

	// ctor is the function that creates an Algorithm under test
	ctor func() *Algorithm[HR]

	// testHash creates an exactly equivalent hash to the ctor's algorithm,
	// to verify behavior of the algorithm.
	testHash func() (expected hash.Hash, expectedSum func() HR)

	// sum is the package-level SumXX function that is equivalent to what ctor returns.
	sum func([]byte) HR
}

func (suite *AlgorithmTestSuite[HR]) SetupTest() {
	suite.testString = "these are some lovely test bytes"
	suite.testBytes = []byte(suite.testString)
	suite.Require().NotNil(suite.ctor)
	suite.Require().NotNil(suite.testHash)
	suite.Require().NotNil(suite.sum)
}

func (suite *AlgorithmTestSuite[HR]) TestNew() {
	actual := suite.ctor().New()
	expected, expectedSum := suite.testHash()

	expected.Write(suite.testBytes)
	actual.Write(suite.testBytes)
	suite.Equal(expectedSum(), actual.Value())
}

func (suite *AlgorithmTestSuite[HR]) TestSum() {
	actual := suite.ctor().Sum(suite.testBytes)
	suite.Equal(suite.sum(suite.testBytes), actual)
}

func (suite *AlgorithmTestSuite[HR]) TestSumString() {
	actual := suite.ctor().SumString(suite.testString)
	expected := suite.sum([]byte(suite.testString))
	suite.Equal(expected, actual)
}

// asSum32 returns a low-level 32-bit Sum function using a hash.Hash32.
// Used for testing to get around: https://github.com/spaolacci/murmur3/issues/34
func asSum32(ctor func() hash.Hash32) func([]byte) uint32 {
	return func(b []byte) uint32 {
		h := ctor()
		h.Write(b)
		return h.Sum32()
	}
}

// asSum64 returns a low-level 32-bit Sum function using a hash.Hash32.
// Used for testing to get around: https://github.com/spaolacci/murmur3/issues/34
func asSum64(ctor func() hash.Hash64) func([]byte) uint64 {
	return func(b []byte) uint64 {
		h := ctor()
		h.Write(b)
		return h.Sum64()
	}
}

func TestNewAlgorithm(t *testing.T) {
	t.Run("NilConstructor", func(t *testing.T) {
		assert.Panics(t, func() {
			NewAlgorithm(nil, murmur3.Sum32)
		})
	})
}

// TestDefault32 runs deeper tests on the Algorithm type as a whole,
// using medley's default 32-bit algorithm.
func TestDefault32(t *testing.T) {
	suite.Run(t, &AlgorithmTestSuite[uint32]{
		ctor: Default32,
		testHash: func() (expected hash.Hash, expectedSum func() uint32) {
			h := murmur3.New32()
			return h, h.Sum32
		},

		// We can't use the murmur3.Sum32 function right now because of:
		// https://github.com/spaolacci/murmur3/issues/34
		sum: asSum32(murmur3.New32),
	})

	t.Run("NilSum", func(t *testing.T) {
		suite.Run(t, &AlgorithmTestSuite[uint32]{
			ctor: func() *Algorithm[uint32] {
				// check that synthesizing a Sum([]byte) HR funtion works as intended
				return NewAlgorithm(AsConstructor32(murmur3.New32), nil)
			},
			testHash: func() (expected hash.Hash, expectedSum func() uint32) {
				h := murmur3.New32()
				return h, h.Sum32
			},

			// We can't use the murmur3.Sum32 function right now because of:
			// https://github.com/spaolacci/murmur3/issues/34
			sum: asSum32(murmur3.New32),
		})
	})
}

// TestDefault32 runs deeper tests on the Algorithm type as a whole,
// using medley's default 32-bit algorithm.
func TestDefault64(t *testing.T) {
	suite.Run(t, &AlgorithmTestSuite[uint64]{
		ctor: Default64,
		testHash: func() (expected hash.Hash, expectedSum func() uint64) {
			h := murmur3.New64()
			return h, h.Sum64
		},

		// We can't use the murmur3.Sum32 function right now because of:
		// https://github.com/spaolacci/murmur3/issues/34
		sum: asSum64(murmur3.New64),
	})

	t.Run("NilSum", func(t *testing.T) {
		suite.Run(t, &AlgorithmTestSuite[uint64]{
			ctor: func() *Algorithm[uint64] {
				// check that synthesizing a Sum([]byte) HR funtion works as intended
				return NewAlgorithm(AsConstructor64(murmur3.New64), nil)
			},
			testHash: func() (expected hash.Hash, expectedSum func() uint64) {
				h := murmur3.New64()
				return h, h.Sum64
			},

			// We can't use the murmur3.Sum32 function right now because of:
			// https://github.com/spaolacci/murmur3/issues/34
			sum: asSum64(murmur3.New64),
		})
	})
}

func TestFNV32(t *testing.T) {
	suite.Run(t, &AlgorithmTestSuite[uint32]{
		ctor: FNV32,
		testHash: func() (expected hash.Hash, expectedSum func() uint32) {
			h := fnv.New32()
			return h, h.Sum32
		},
		// the hash/fnv package does not provide a package-level sum function
		sum: asSum32(fnv.New32),
	})
}

func TestFNV32a(t *testing.T) {
	suite.Run(t, &AlgorithmTestSuite[uint32]{
		ctor: FNV32a,
		testHash: func() (expected hash.Hash, expectedSum func() uint32) {
			h := fnv.New32a()
			return h, h.Sum32
		},
		// the hash/fnv package does not provide a package-level sum function
		sum: asSum32(fnv.New32a),
	})
}

func TestFNV64(t *testing.T) {
	suite.Run(t, &AlgorithmTestSuite[uint64]{
		ctor: FNV64,
		testHash: func() (expected hash.Hash, expectedSum func() uint64) {
			h := fnv.New64()
			return h, h.Sum64
		},
		// the hash/fnv package does not provide a package-level sum function
		sum: asSum64(fnv.New64),
	})
}

func TestFNV64a(t *testing.T) {
	suite.Run(t, &AlgorithmTestSuite[uint64]{
		ctor: FNV64a,
		testHash: func() (expected hash.Hash, expectedSum func() uint64) {
			h := fnv.New64a()
			return h, h.Sum64
		},
		// the hash/fnv package does not provide a package-level sum function
		sum: asSum64(fnv.New64a),
	})
}
