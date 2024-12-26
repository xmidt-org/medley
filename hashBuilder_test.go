// SPDX-FileCopyrightText: 2025 Comcast Cable Communications Management, LLC
// SPDX-License-Identifier: Apache-2.0

package medley

import (
	"bytes"
	"encoding/binary"
	"hash/fnv"
	"io"
	"math"
	"testing"

	"github.com/stretchr/testify/suite"
)

type HashBuilderSuite struct {
	suite.Suite
}

func (suite *HashBuilderSuite) newHashBuilder(dst io.Writer) *HashBuilder {
	hb := NewHashBuilder(dst)
	suite.Require().NotNil(hb)
	return hb
}

func (suite *HashBuilderSuite) assertWriteSuccess(hb *HashBuilder) *HashBuilder {
	suite.Require().NotNil(hb)
	suite.Require().NoError(hb.Err())
	return hb
}

func (suite *HashBuilderSuite) TestWrite() {
	var b bytes.Buffer
	hb := suite.newHashBuilder(&b)

	hb = suite.assertWriteSuccess(
		hb.Write([]byte{'a', 'b', 'c', 'd'}),
	)

	suite.Equal(
		[]byte{'a', 'b', 'c', 'd'},
		b.Bytes(),
	)

	hb.Reset()
	suite.Empty(b.Bytes())
}

func (suite *HashBuilderSuite) TestWriteString() {
	var b bytes.Buffer
	hb := suite.newHashBuilder(&b)

	hb = suite.assertWriteSuccess(
		hb.WriteString("test value"),
	)

	suite.Equal("test value", b.String())
	hb.Reset()
	suite.Empty(b.Bytes())
}

func (suite *HashBuilderSuite) TestWriteUint8() {
	const expected uint8 = '0'
	var b bytes.Buffer
	hb := suite.newHashBuilder(&b)

	hb = suite.assertWriteSuccess(
		hb.WriteUint8(expected),
	)

	suite.Equal([]byte{expected}, b.Bytes())
	hb.Reset()
	suite.Empty(b.Bytes())
}

func (suite *HashBuilderSuite) TestWriteUint16() {
	const expected uint16 = 457
	var b bytes.Buffer
	hb := suite.newHashBuilder(&b)

	hb = suite.assertWriteSuccess(
		hb.WriteUint16(expected),
	)

	suite.Equal(
		expected,
		binary.BigEndian.Uint16(b.Bytes()),
	)

	hb.Reset()
	suite.Empty(b.Bytes())
}

func (suite *HashBuilderSuite) TestWriteUint32() {
	const expected uint32 = 34987342
	var b bytes.Buffer
	hb := suite.newHashBuilder(&b)

	hb = suite.assertWriteSuccess(
		hb.WriteUint32(expected),
	)

	suite.Equal(
		expected,
		binary.BigEndian.Uint32(b.Bytes()),
	)

	hb.Reset()
	suite.Empty(b.Bytes())
}

func (suite *HashBuilderSuite) TestWriteUint64() {
	const expected uint64 = 2957103845918235
	var b bytes.Buffer
	hb := suite.newHashBuilder(&b)

	hb = suite.assertWriteSuccess(
		hb.WriteUint64(expected),
	)

	suite.Equal(
		expected,
		binary.BigEndian.Uint64(b.Bytes()),
	)

	hb.Reset()
	suite.Empty(b.Bytes())
}

func (suite *HashBuilderSuite) TestWriteFloat32() {
	const expected float32 = -523974.498234
	var b bytes.Buffer
	hb := suite.newHashBuilder(&b)

	hb = suite.assertWriteSuccess(
		hb.WriteFloat32(expected),
	)

	suite.Equal(
		expected,
		math.Float32frombits(
			binary.BigEndian.Uint32(b.Bytes()),
		),
	)

	hb.Reset()
	suite.Empty(b.Bytes())
}

func (suite *HashBuilderSuite) TestWriteFloat64() {
	const expected float64 = 234984534535272336.3497834534345
	var b bytes.Buffer
	hb := suite.newHashBuilder(&b)

	hb = suite.assertWriteSuccess(
		hb.WriteFloat64(expected),
	)

	suite.Equal(
		expected,
		math.Float64frombits(
			binary.BigEndian.Uint64(b.Bytes()),
		),
	)

	hb.Reset()
	suite.Empty(b.Bytes())
}

func (suite *HashBuilderSuite) writeHashBytes(hb *HashBuilder) int {
	suite.Require().NotNil(hb)
	hb.
		Write([]byte{'1', '2', '3'}). // 3 bytes
		WriteString("test").          // 4 bytes (7 total)
		WriteUint8(123).              // 1 byte (8 total)
		WriteUint16(32948).           // 2 bytes (10 total)
		WriteUint32(348901223).       // 4 bytes (14 total)
		WriteUint64(23394580232312).  // 8 bytes (22 total)
		WriteFloat32(47.6).           // 4 bytes (26 total)
		WriteFloat64(-2342135.23423)  // 8 bytes (34 total)

	suite.Require().NoError(hb.Err())

	return 34 // total written
}

func (suite *HashBuilderSuite) TestWithBuffer() {
	var b bytes.Buffer
	hb := NewHashBuilder(&b)
	suite.True(hb.CanReset())
	suite.False(hb.CanSum64())

	suite.Zero(hb.Sum64())
	expectedLength := suite.writeHashBytes(hb)
	suite.Equal(expectedLength, b.Len())
	suite.Zero(hb.Sum64())

	hb.Reset()
	suite.Zero(b.Len())
	suite.Zero(hb.Sum64())
}

func (suite *HashBuilderSuite) TestWithHash() {
	hash := fnv.New64()
	initial := hash.Sum64()

	hb := NewHashBuilder(hash)
	suite.True(hb.CanReset())
	suite.True(hb.CanSum64())

	suite.Equal(initial, hb.Sum64())
	suite.Equal(hb.Sum64(), hash.Sum64())
	suite.writeHashBytes(hb)
	suite.NotZero(hb.Sum64())
	suite.Equal(hb.Sum64(), hash.Sum64())

	hb.Reset()
	suite.Equal(initial, hb.Sum64())
	suite.Equal(hb.Sum64(), hash.Sum64())
}

func TestHashBuilder(t *testing.T) {
	suite.Run(t, new(HashBuilderSuite))
}
