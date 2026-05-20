// SPDX-FileCopyrightText: 2026 Comcast Cable Communications Management, LLC
// SPDX-License-Identifier: Apache-2.0

package medley

import (
	"hash"
	"hash/crc32"
	"hash/crc64"
	"hash/fnv"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// TestHashValuesSuite holds common test values for hashing.
type TestHashValuesSuite struct {
	suite.Suite

	testBytes  []byte
	testString string
}

func (suite *TestHashValuesSuite) SetupSuite() {
	suite.testString = "here are some lovely test bytes"
	suite.testBytes = []byte(suite.testString)
}

// assertWriteTestByte writes the first byte of suite.testBytes to dst and verifies the result.
func (suite *TestHashValuesSuite) assertWriteTestByte(dst io.ByteWriter) {
	suite.Require().Greater(len(suite.testBytes), 0) // guard against misconfiguration of the suite

	suite.NoError(
		dst.WriteByte(suite.testBytes[0]),
	)
}

// assertWriteTestBytes writes the test bytes to dst and verifies the result.
func (suite *TestHashValuesSuite) assertWriteTestBytes(dst io.Writer) {
	suite.Require().Greater(len(suite.testBytes), 0) // guard against misconfiguration of the suite

	n, err := dst.Write(suite.testBytes)
	suite.Equal(len(suite.testBytes), n)
	suite.NoError(err)
}

// assertWriteTestBytes writes the test string to dst and verifies the result.
func (suite *TestHashValuesSuite) assertWriteTestString(dst io.StringWriter) {
	suite.Require().Greater(len(suite.testString), 0) // guard against misconfiguration of the suite

	n, err := dst.WriteString(suite.testString)
	suite.Equal(len(suite.testString), n)
	suite.NoError(err)
}

// HashTestSuite tests the Hash[HR] type specifically.
type HashTestSuite[HR HashResult] struct {
	TestHashValuesSuite

	// ctor is the appropriate ConstructorXX function
	ctor func() Hash[HR]

	// testHash creates a hash under test.  expected and expectedSum are the hash and SumXXX()
	// methods of the underlying hash.Hash, and actual is the medley Hash object that wraps expected.
	testHash func() (expected hash.Hash, expectedSum func() HR, actual Hash[HR])
}

func (suite *HashTestSuite[HR]) SetupSuite() {
	suite.TestHashValuesSuite.SetupSuite()
	suite.Require().NotNil(suite.ctor)
	suite.Require().NotNil(suite.testHash)
}

// newTestHash creates a hash.Hash with an actual Hash[HR] that wraps it. Basic assertions
// are done to make sure Hash[HR] does not modify the hash.
func (suite *HashTestSuite[HR]) newTestHash() (expected hash.Hash, expectedSum func() HR, actual Hash[HR]) {
	expected, expectedSum, actual = suite.testHash()
	suite.Require().Equal(expected.BlockSize(), actual.BlockSize())
	suite.Require().Equal(expected.Size(), actual.Size())

	return
}

func (suite *HashTestSuite[HR]) TestConstructor() {
	expected, expectedSum, _ := suite.newTestHash()
	actual := suite.ctor()
	suite.Require().Equal(expected.BlockSize(), actual.BlockSize())
	suite.Require().Equal(expected.Size(), actual.Size())

	suite.assertWriteTestBytes(expected)
	suite.assertWriteTestBytes(actual)
	suite.Equal(expectedSum(), actual.Value())
}

func (suite *HashTestSuite[HR]) TestSum() {
	expected, _, actual := suite.newTestHash()
	expectedBytes := expected.Sum(suite.testBytes)
	actualBytes := actual.Sum(suite.testBytes)
	suite.Equal(expectedBytes, actualBytes)
}

func (suite *HashTestSuite[HR]) TestWrite() {
	expected, expectedSum, actual := suite.newTestHash()
	suite.assertWriteTestBytes(expected)
	suite.Equal(expectedSum(), actual.Value())

	// verify that Reset works correctly
	expected.Reset()
	suite.assertWriteTestBytes(expected)
	suite.Equal(expectedSum(), actual.Value())
}

func (suite *HashTestSuite[HR]) TestWriteString() {
	_, expectedSum, actual := suite.newTestHash()
	initial := expectedSum()

	suite.assertWriteTestString(actual)
	suite.NotEqual(initial, expectedSum())
	suite.Equal(expectedSum(), actual.Value())

	// verify that Reset works correctly
	actual.Reset()
	suite.Equal(initial, expectedSum())
	suite.assertWriteTestString(actual)
	suite.NotEqual(initial, expectedSum())
	suite.Equal(expectedSum(), actual.Value())
}

func (suite *HashTestSuite[HR]) TestWriteByte() {
	_, expectedSum, actual := suite.newTestHash()
	initial := expectedSum()

	suite.assertWriteTestByte(actual)
	suite.NotEqual(initial, expectedSum())
	suite.Equal(expectedSum(), actual.Value())

	// verify that Reset works correctly
	actual.Reset()
	suite.Equal(initial, expectedSum())
	suite.assertWriteTestByte(actual)
	suite.NotEqual(initial, expectedSum())
	suite.Equal(expectedSum(), actual.Value())
}

func TestHash32(t *testing.T) {
	suite.Run(t, &HashTestSuite[uint32]{
		ctor: AsConstructor32(crc32.NewIEEE),
		testHash: func() (hash.Hash, func() uint32, Hash[uint32]) {
			h := crc32.NewIEEE()
			return h, h.Sum32, AsHash32(h)
		},
	})

	t.Run("AsContructorNil", func(t *testing.T) {
		assert.Panics(t, func() {
			AsConstructor32(nil)
		})
	})
}

func TestHash64(t *testing.T) {
	table := crc64.MakeTable(0x23FEAD) // just a random CRC table

	suite.Run(t, &HashTestSuite[uint64]{
		ctor: AsConstructor64(func() hash.Hash64 { return crc64.New(table) }),
		testHash: func() (hash.Hash, func() uint64, Hash[uint64]) {
			h := crc64.New(table)
			return h, h.Sum64, AsHash64(h)
		},
	})

	t.Run("AsContructorNil", func(t *testing.T) {
		assert.Panics(t, func() {
			AsConstructor64(nil)
		})
	})
}

type SumTestSuite[HR HashResult] struct {
	TestHashValuesSuite

	ctor func() Hash[HR]
}

func (suite *SumTestSuite[HR]) SetupSuite() {
	suite.TestHashValuesSuite.SetupSuite()
	suite.Require().NotNil(suite.ctor)
}

func (suite *SumTestSuite[HR]) TestAsSum() {
	sum := AsSum(suite.ctor)
	suite.Require().NotNil(sum)

	h := suite.ctor()
	suite.Require().NotNil(h)
	suite.assertWriteTestBytes(h)

	suite.Equal(h.Value(), sum(suite.testBytes))
}

func (suite *SumTestSuite[HR]) TestSumString() {
	sum := AsSum(suite.ctor)
	suite.Require().NotNil(sum)

	h := suite.ctor()
	suite.Require().NotNil(h)
	suite.assertWriteTestString(h)

	suite.Equal(h.Value(), SumString(sum, suite.testString))
}

func TestSum32(t *testing.T) {
	suite.Run(t, &SumTestSuite[uint32]{
		ctor: AsConstructor32(fnv.New32a),
	})
}

func TestSum64(t *testing.T) {
	suite.Run(t, &SumTestSuite[uint64]{
		ctor: AsConstructor64(fnv.New64a),
	})
}

// AlgorithmTestSuite tests the builtin medley algorithm functions that produce
// hashing objects, e.g. Default64().
type AlgorithmTestSuite[HR HashResult] struct {
	TestHashValuesSuite

	alg func() (Constructor[HR], Sum[HR])
}

func (suite *AlgorithmTestSuite[HR]) SetupSuite() {
	suite.TestHashValuesSuite.SetupSuite()
	suite.Require().NotNil(suite.alg)
}

func (suite *AlgorithmTestSuite[HR]) TestBytes() {
	ctor, sum := suite.alg()
	suite.Require().NotNil(ctor)
	suite.Require().NotNil(sum)

	h := ctor()
	suite.assertWriteTestBytes(h)

	suite.Equal(h.Value(), sum(suite.testBytes))
}

func (suite *AlgorithmTestSuite[HR]) TestString() {
	ctor, sum := suite.alg()
	suite.Require().NotNil(ctor)
	suite.Require().NotNil(sum)

	h := ctor()
	suite.assertWriteTestString(h)

	suite.Equal(h.Value(), SumString(sum, suite.testString))
}

func TestAlgorithm(t *testing.T) {
	t.Run("Default32", func(t *testing.T) {
		suite.Run(t, &AlgorithmTestSuite[uint32]{
			alg: Default32,
		})
	})

	t.Run("Default64", func(t *testing.T) {
		suite.Run(t, &AlgorithmTestSuite[uint64]{
			alg: Default64,
		})
	})

	t.Run("FNV32a", func(t *testing.T) {
		suite.Run(t, &AlgorithmTestSuite[uint32]{
			alg: FNV32a,
		})
	})

	t.Run("FNV64a", func(t *testing.T) {
		suite.Run(t, &AlgorithmTestSuite[uint64]{
			alg: FNV64a,
		})
	})
}
