// SPDX-FileCopyrightText: 2026 Comcast Cable Communications Management, LLC
// SPDX-License-Identifier: Apache-2.0

package medley

import (
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func TestAppend(t *testing.T) {
	const testValue = "here is a test value"

	t.Run("ByteSlice", func(t *testing.T) {
		object := []byte(testValue)
		var dst []byte
		dst = Append(dst, object)
		assert.Equal(t, object, dst)
	})

	t.Run("CustomByteSlice", func(t *testing.T) {
		type Custom []byte
		object := Custom(testValue)
		var dst []byte
		dst = Append(dst, object)
		assert.Equal(t, []byte(object), dst)
	})

	t.Run("String", func(t *testing.T) {
		object := testValue
		var dst []byte
		dst = Append(dst, object)
		assert.Equal(t, []byte(object), dst)
	})

	t.Run("CustomString", func(t *testing.T) {
		type Custom string
		object := Custom(testValue)
		var dst []byte
		dst = Append(dst, object)
		assert.Equal(t, []byte(object), dst)
	})
}

// ObjectifyTestSuite tests the generic Objectify function.
type ObjectifyTestSuite struct {
	suite.Suite
}

func (suite *ObjectifyTestSuite) TestSequence() {
	type server struct {
		hostName string
		port     int
	}

	objecter := func(s *server) string { return s.hostName }

	var servers = []*server{
		{"host1.something.net", 1111},
		{"host2.something.net", 2222},
		{"host3.something.net", 3333},
	}

	suite.Run("Iteration", func() {
		var i int
		for object, value := range Objectify(objecter, slices.Values(servers)) {
			suite.Require().Less(i, len(servers))
			suite.Equal(servers[i].hostName, object)
			suite.Equal(servers[i].port, value.port)
			suite.Equal(servers[i], value)
			i++
		}

		suite.Equal(len(servers), i)
	})

	suite.Run("EarlyReturn", func() {
		var i int
		for object, value := range Objectify(objecter, slices.Values(servers)) {
			suite.Require().Zero(i)
			suite.Equal(servers[i].hostName, object)
			suite.Equal(servers[i], value)
			i++
			break
		}

		suite.Equal(1, i)
	})
}

func TestObjectify(t *testing.T) {
	suite.Run(t, new(ObjectifyTestSuite))
}

// StringifyTestSuite tests the Stringify function.
type StringifyTestSuite struct {
	suite.Suite
}

func (suite *StringifyTestSuite) TestSequence() {
	var hostNames = []string{
		"host1.something.net",
		"host2.something.net",
		"host3.something.net",
	}

	suite.Run("Iteration", func() {
		var i int
		for object, value := range Stringify(slices.Values(hostNames)) {
			suite.Require().Less(i, len(hostNames))
			suite.Equal(hostNames[i], object)
			suite.Equal(hostNames[i], value)
			i++
		}

		suite.Equal(len(hostNames), i)
	})

	suite.Run("EarlyReturn", func() {
		var i int
		for object, value := range Stringify(slices.Values(hostNames)) {
			suite.Require().Zero(i)
			suite.Equal(hostNames[i], object)
			suite.Equal(hostNames[i], value)
			i++
			break
		}

		suite.Equal(1, i)
	})
}

func TestStringify(t *testing.T) {
	suite.Run(t, new(StringifyTestSuite))
}
