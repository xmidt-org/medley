// SPDX-FileCopyrightText: 2026 Comcast Cable Communications Management, LLC
// SPDX-License-Identifier: Apache-2.0

package medley

import (
	"bytes"
	"encoding/binary"
	"iter"
	"reflect"
	"slices"
	"strconv"
	"testing"

	"github.com/stretchr/testify/suite"
)

// ObjectTestSuite holds common infrastructure for testing Objects of
// any kind.
type ObjectTestSuite[C any] struct {
	suite.Suite
}

// assertLen verifies that Object.Len behaves correctly.
func (suite *ObjectTestSuite[C]) assertLen(expectedLen int, obj Object) {
	suite.Equal(expectedLen, obj.Len())
	buf := obj.Append([]byte{})
	suite.Equal(expectedLen, len(buf))
	suite.Equal(obj.b, buf)
}

// assertToHash32 verifies that a lifecycle involving the object's ToHash
// works correctly with a Hash32.
func (suite *ObjectTestSuite[C]) assertToHash(obj Object) {
	var buffer bytes.Buffer
	obj.ToHash(&buffer)
	if obj.Len() == 0 {
		suite.Zero(buffer.Len())
	} else {
		suite.Equal(obj.b, buffer.Bytes())
	}
}

// assertWriteTo verifies that WriterTo behaves correct for the given object.
func (suite *ObjectTestSuite[C]) assertWriteTo(obj Object) {
	var buffer bytes.Buffer
	n, err := obj.WriteTo(&buffer)
	suite.Equal(obj.Len(), int(n))
	suite.NoError(err)
}

type arbitraryLengthTestCase[C []byte | string] struct {
	name     string
	contents C
}

// ArbitraryLengthObjectTestSuite holds common infrastructure for Objects whose
// length can vary.
type ArbitraryLengthObjectTestSuite[C []byte | string] struct {
	ObjectTestSuite[C]

	testCases []arbitraryLengthTestCase[C]
}

type BytesTestSuite struct {
	ArbitraryLengthObjectTestSuite[[]byte]
}

func (suite *BytesTestSuite) SetupSuite() {
	suite.testCases = []arbitraryLengthTestCase[[]byte]{
		{
			name: "nil",
		},
		{
			name:     "empty",
			contents: []byte{},
		},
		{
			name:     "1",
			contents: []byte{123},
		},
		{
			name:     "5",
			contents: []byte{12, 78, 191, 45, 254},
		},
		{
			name: "20",
			contents: []byte{
				56, 143, 90, 178, 1,
				67, 23, 83, 217, 198,
				194, 4, 17, 54, 32,
				235, 209, 11, 78, 176,
			},
		},
	}
}

func (suite *BytesTestSuite) TestLen() {
	for _, testCase := range suite.testCases {
		suite.Run(testCase.name, func() {
			obj := Bytes(testCase.contents)
			suite.assertLen(len(testCase.contents), obj)
		})
	}
}

func (suite *BytesTestSuite) TestToHash() {
	for _, testCase := range suite.testCases {
		suite.Run(testCase.name, func() {
			obj := Bytes(testCase.contents)
			suite.assertToHash(obj)
		})
	}
}

func (suite *BytesTestSuite) TestWriteTo() {
	for _, testCase := range suite.testCases {
		suite.Run(testCase.name, func() {
			obj := Bytes(testCase.contents)
			suite.assertWriteTo(obj)
		})
	}
}

func TestBytes(t *testing.T) {
	suite.Run(t, new(BytesTestSuite))
}

type StringTestSuite struct {
	ArbitraryLengthObjectTestSuite[string]
}

func (suite *StringTestSuite) SetupSuite() {
	suite.testCases = []arbitraryLengthTestCase[string]{
		{
			name: "uninitialized",
		},
		{
			name:     "empty",
			contents: "",
		},
		{
			name:     "1",
			contents: "a",
		},
		{
			name:     strconv.Itoa(len("chair")),
			contents: "chair",
		},
		{
			name:     strconv.Itoa(len("the quick brown fox")),
			contents: "the quick brown fox",
		},
	}
}

func (suite *StringTestSuite) TestLen() {
	for _, testCase := range suite.testCases {
		suite.Run(testCase.name, func() {
			obj := String(testCase.contents)
			suite.assertLen(len(testCase.contents), obj)
		})
	}
}

func (suite *StringTestSuite) TestToHash() {
	for _, testCase := range suite.testCases {
		suite.Run(testCase.name, func() {
			obj := String(testCase.contents)
			suite.assertToHash(obj)
		})
	}
}

func (suite *StringTestSuite) TestWriteTo() {
	for _, testCase := range suite.testCases {
		suite.Run(testCase.name, func() {
			obj := String(testCase.contents)
			suite.assertWriteTo(obj)
		})
	}
}

func TestString(t *testing.T) {
	suite.Run(t, new(StringTestSuite))
}

type integerTestCase[U uint16 | uint32 | uint64] struct {
	name      string
	contents  U
	byteOrder binary.ByteOrder
}

// IntegerTestSuite runs tests over the Integer constructor for objects.
// Input is fixed length, making much fewer test cases.
type IntegerTestSuite[U uint16 | uint32 | uint64] struct {
	ObjectTestSuite[U]

	expectedLen int
	testCases   []integerTestCase[U]
}

func (suite *IntegerTestSuite[U]) SetupSuite() {
	suite.expectedLen = int(reflect.TypeFor[U]().Size())

	// testBytes is just a constant set of bytes used
	// to generate uints of the various sizes we support.
	testBytes := [8]byte{
		0xF5, 0x39, 0xAE, 0x19,
		0xD4, 0x5B, 0x95, 0xDC,
	}

	var testValue U
	for i := range suite.expectedLen {
		testValue <<= 8
		testValue |= U(testBytes[i])
	}

	suite.testCases = []integerTestCase[U]{
		{
			name:      "BigEndian",
			contents:  testValue,
			byteOrder: binary.BigEndian,
		},
		{
			name:      "LittleEndian",
			contents:  testValue,
			byteOrder: binary.LittleEndian,
		},
	}
}

func (suite *IntegerTestSuite[U]) TestLen() {
	for _, testCase := range suite.testCases {
		suite.Run(testCase.name, func() {
			obj := Integer(testCase.contents, testCase.byteOrder)
			suite.assertLen(suite.expectedLen, obj)
		})
	}
}

func (suite *IntegerTestSuite[U]) TestToHash() {
	for _, testCase := range suite.testCases {
		suite.Run(testCase.name, func() {
			obj := Integer(testCase.contents, testCase.byteOrder)
			suite.assertToHash(obj)
		})
	}
}

func (suite *IntegerTestSuite[U]) TestWriteTo() {
	for _, testCase := range suite.testCases {
		suite.Run(testCase.name, func() {
			obj := Integer(testCase.contents, testCase.byteOrder)
			suite.assertWriteTo(obj)
		})
	}
}

func TestInteger16(t *testing.T) {
	suite.Run(t, new(IntegerTestSuite[uint16]))
}

func TestInteger32(t *testing.T) {
	suite.Run(t, new(IntegerTestSuite[uint32]))
}

func TestInteger64(t *testing.T) {
	suite.Run(t, new(IntegerTestSuite[uint64]))
}

// ObjectSequenceTestSuite holds common infrastructure for testing sequences of
// hashable objects produced by Objectify and its variants.
type ObjectSequenceTestSuite[V comparable] struct {
	suite.Suite
}

// assertSequence verifies that the actual Objectify-style sequence visits each of the expected
// elements in proper order.
func (suite *ObjectSequenceTestSuite[V]) assertSequence(expectedCount int, verify Objecter[V], expected iter.Seq[V], actual iter.Seq2[Object, V]) {
	expectedNext, expectedStop := iter.Pull(expected)
	defer expectedStop()

	actualNext, actualStop := iter.Pull2(actual)
	defer actualStop()

	actualCount := 0
	for expected, ok := expectedNext(); ok; expected, ok = expectedNext() {
		actualObject, actual, ok := actualNext()
		suite.Require().True(ok)

		actualCount++
		suite.Equal(expected, actual)
		suite.Equal(
			verify(actual).b,
			actualObject.b,
		)
	}

	suite.Equal(expectedCount, actualCount)
}

// assertSlice verifies that the actual Objectify-style sequence visits each of the expected
// slice elements in proper order.
func (suite *ObjectSequenceTestSuite[V]) assertSlice(verify Objecter[V], expected []V, actual iter.Seq2[Object, V]) {
	actualNext, actualStop := iter.Pull2(actual)
	defer actualStop()

	for i := 0; i < len(expected); i++ {
		expected := expected[i]

		actualObject, actual, ok := actualNext()
		suite.Require().True(ok)

		suite.Equal(expected, actual)
		suite.Equal(
			verify(actual).b,
			actualObject.b,
		)
	}
}

// assertEarlyReturn verifies that the yield function works correctly and will
// allow a for loop to break early.
func (suite *ObjectSequenceTestSuite[V]) assertEarlyReturn(expectedFirstObject Object, expectedFirstValue V, actual iter.Seq2[Object, V]) {
	iterations := 0
	for obj, value := range actual {
		suite.Require().Zero(iterations)
		suite.Equal(expectedFirstObject.b, obj.b)
		suite.Equal(expectedFirstValue, value)
		iterations++
		break
	}

	suite.Equal(iterations, 1)
}

// ObjectifyTestSuite tests the generic Objectify and ObjectifySlice functions.
type ObjectifyTestSuite struct {
	ObjectSequenceTestSuite[string]
}

func (suite *ObjectifyTestSuite) TestSequence() {
	testValues := slices.Values([]string{"one", "two", "three"})

	suite.assertSequence(
		3,
		String,
		testValues,
		Objectify(
			String,
			testValues,
		),
	)
}

func (suite *ObjectifyTestSuite) TestSlice() {
	testValues := []string{"one", "two", "three"}

	suite.assertSlice(
		String,
		testValues,
		ObjectifySlice(
			String,
			testValues,
		),
	)
}

func (suite *ObjectifyTestSuite) TestEarlyReturn() {
	suite.Run("Sequence", func() {
		testValues := []string{"one", "two", "three"}
		suite.assertEarlyReturn(
			String("one"),
			"one",
			Objectify(String, slices.Values(testValues)),
		)
	})

	suite.Run("Slice", func() {
		testValues := []string{"one", "two", "three"}
		suite.assertEarlyReturn(
			String("one"),
			"one",
			ObjectifySlice(String, testValues),
		)
	})
}

func TestObjectify(t *testing.T) {
	suite.Run(t, new(ObjectifyTestSuite))
}

type StringifyTestSuite struct {
	ObjectSequenceTestSuite[string]
}

func (suite *StringifyTestSuite) TestSequence() {
	testValues := slices.Values([]string{"one", "two", "three"})

	suite.assertSequence(
		3,
		String,
		testValues,
		Stringify(testValues),
	)
}

func (suite *StringifyTestSuite) TestSlice() {
	testValues := []string{"one", "two", "three"}

	suite.assertSlice(
		String,
		testValues,
		StringifySlice(testValues),
	)
}

func (suite *StringifyTestSuite) TestEarlyReturn() {
	suite.Run("Sequence", func() {
		testValues := []string{"one", "two", "three"}
		suite.assertEarlyReturn(
			String("one"),
			"one",
			Stringify(slices.Values(testValues)),
		)
	})

	suite.Run("Slice", func() {
		testValues := []string{"one", "two", "three"}
		suite.assertEarlyReturn(
			String("one"),
			"one",
			StringifySlice(testValues),
		)
	})
}

func TestStringify(t *testing.T) {
	suite.Run(t, new(StringifyTestSuite))
}
