// SPDX-FileCopyrightText: 2026 Comcast Cable Communications Management, LLC
// SPDX-License-Identifier: Apache-2.0

package medley

import (
	"hash"
	"testing"

	"github.com/spaolacci/murmur3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type AlgorithmSuite[HR HashResult] struct {
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

func (suite *AlgorithmSuite[HR]) SetupTest() {
	suite.testString = "these are some lovely test bytes"
	suite.testBytes = []byte(suite.testString)
	suite.Require().NotNil(suite.ctor)
	suite.Require().NotNil(suite.testHash)
	suite.Require().NotNil(suite.sum)
}

func (suite *AlgorithmSuite[HR]) TestNew() {
	actual := suite.ctor().New()
	expected, expectedSum := suite.testHash()

	expected.Write(suite.testBytes)
	actual.Write(suite.testBytes)
	suite.Equal(expectedSum(), actual.Value())
}

func (suite *AlgorithmSuite[HR]) TestSum() {
	actual := suite.ctor().Sum(suite.testBytes)
	suite.Equal(suite.sum(suite.testBytes), actual)
}

func (suite *AlgorithmSuite[HR]) TestSumString() {
	actual := suite.ctor().SumString(suite.testString)
	expected := suite.sum([]byte(suite.testString))
	suite.Equal(expected, actual)
}

func (suite *AlgorithmSuite[HR]) TestSumObject() {
	obj := Bytes(suite.testBytes)
	actual := suite.ctor().SumObject(obj)
	expected := suite.sum(obj.b)
	suite.Equal(expected, actual)
}

func TestAlgorithm32(t *testing.T) {
	t.Run("Default32", func(t *testing.T) {
		suite.Run(t, &AlgorithmSuite[uint32]{
			ctor: Default32,
			testHash: func() (expected hash.Hash, expectedSum func() uint32) {
				h := murmur3.New32()
				return h, h.Sum32
			},
			sum: murmur3.Sum32,
		})
	})

	t.Run("NilSum", func(t *testing.T) {
		suite.Run(t, &AlgorithmSuite[uint32]{
			ctor: func() *Algorithm[uint32] {
				// check that synthesizing a Sum([]byte) HR funtion works as intended
				return NewAlgorithm(AsConstructor32(murmur3.New32), nil)
			},
			testHash: func() (expected hash.Hash, expectedSum func() uint32) {
				h := murmur3.New32()
				return h, h.Sum32
			},
			sum: murmur3.Sum32,
		})
	})

	t.Run("NewAlgorithmNilConstructor", func(t *testing.T) {
		assert.Panics(t, func() {
			NewAlgorithm[uint32](nil, murmur3.Sum32)
		})
	})
}

func TestAlgorithm64(t *testing.T) {
	t.Run("Default32", func(t *testing.T) {
		suite.Run(t, &AlgorithmSuite[uint64]{
			ctor: Default64,
			testHash: func() (expected hash.Hash, expectedSum func() uint64) {
				h := murmur3.New64()
				return h, h.Sum64
			},
			sum: murmur3.Sum64,
		})
	})

	t.Run("NilSum", func(t *testing.T) {
		suite.Run(t, &AlgorithmSuite[uint64]{
			ctor: func() *Algorithm[uint64] {
				// check that synthesizing a Sum([]byte) HR funtion works as intended
				return NewAlgorithm(AsConstructor64(murmur3.New64), nil)
			},
			testHash: func() (expected hash.Hash, expectedSum func() uint64) {
				h := murmur3.New64()
				return h, h.Sum64
			},
			sum: murmur3.Sum64,
		})
	})

	t.Run("NewAlgorithmNilConstructor", func(t *testing.T) {
		assert.Panics(t, func() {
			NewAlgorithm[uint64](nil, murmur3.Sum64)
		})
	})
}
