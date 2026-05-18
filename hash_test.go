// SPDX-FileCopyrightText: 2026 Comcast Cable Communications Management, LLC
// SPDX-License-Identifier: Apache-2.0

package medley

import (
	"hash"
	"hash/crc32"
	"hash/crc64"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type HashTestSuite[HR HashResult] struct {
	suite.Suite

	testBytes  []byte
	testString string

	// ctor is the appropriate AsConstructorXX function
	ctor func() Hash[HR]

	// testHash creates a hash under test.  expected and expectedSum are the hash and SumXXX()
	// methods of the underlying hash.Hash, and actual is the medley Hash object that wraps expected.
	testHash func() (expected hash.Hash, expectedSum func() HR, actual Hash[HR])
}

func (suite *HashTestSuite[HR]) SetupTest() {
	suite.testString = "here are some test bytes"
	suite.testBytes = []byte(suite.testString)
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

	expected.Write(suite.testBytes)
	actual.Write(suite.testBytes)
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
	expected.Write(suite.testBytes)
	suite.Equal(expectedSum(), actual.Value())

	expected.Reset()
	actual.Write(suite.testBytes)
	suite.Equal(expectedSum(), actual.Value())
}

func (suite *HashTestSuite[HR]) TestWriteString() {
	_, expectedSum, actual := suite.newTestHash()
	initial := expectedSum()

	actual.WriteString(suite.testString)
	suite.NotEqual(initial, expectedSum())
	suite.Equal(expectedSum(), actual.Value())
}

func (suite *HashTestSuite[HR]) TestWriteByte() {
	_, expectedSum, actual := suite.newTestHash()
	initial := expectedSum()

	actual.WriteByte(suite.testBytes[0])
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
