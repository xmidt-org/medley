package consistent

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/xmidt-org/medley"
)

type RingTestSuite struct {
	suite.Suite
}

func (suite *RingTestSuite) TestGrowThenAdd() {
	var ring ring
	suite.Len(ring, 0)
	suite.Equal(0, cap(ring))

	ring.grow(5)
	suite.Len(ring, 0)
	suite.Equal(5, cap(ring))

	ring.add("test", 10)
	suite.Len(ring, 1)
	suite.Equal(5, cap(ring))

	ring.grow(5)
	suite.Len(ring, 1)
	suite.GreaterOrEqual(cap(ring), 10)

	ringCap := cap(ring)
	ring.add("test", 20)
	suite.Len(ring, 2)
	suite.Equal(ringCap, cap(ring))
}

func (suite *RingTestSuite) TestRemoveIf() {
	testData := []struct {
		ring           ring
		p              func(medley.Node) bool
		expectsRemoved int
	}{
		{
			ring:           nil,
			p:              func(medley.Node) bool { return true },
			expectsRemoved: 0,
		},
		{
			ring: ring{
				{Node: "test1", Value: 10}, {Node: "test2", Value: 20}, {Node: "test3", Value: 30}, {Node: "test4", Value: 40},
			},
			p:              func(n medley.Node) bool { return n == "test2" },
			expectsRemoved: 1,
		},
		{
			ring: ring{
				{Node: "test1", Value: 10}, {Node: "test1", Value: 20}, {Node: "test2", Value: 30}, {Node: "test3", Value: 40}, {Node: "test1", Value: 40},
			},
			p:              func(n medley.Node) bool { return n == "test1" },
			expectsRemoved: 3,
		},
		{
			ring: ring{
				{Node: "test1", Value: 10}, {Node: "test2", Value: 20}, {Node: "test3", Value: 30}, {Node: "test4", Value: 40},
			},
			p:              func(n medley.Node) bool { return false },
			expectsRemoved: 0,
		},
	}

	for i, record := range testData {
		suite.Run(strconv.Itoa(i), func() {
			initialLen := record.ring.Len()
			record.ring.removeIf(record.p)
			suite.Len(record.ring, initialLen-record.expectsRemoved)
		})
	}
}

func (suite *RingTestSuite) TestClosest() {
	suite.Run("Empty", func() {
		n, err := ring{}.closest(0)
		suite.Empty(n)
		suite.Error(err)
	})

	testData := []struct {
		ring     ring
		value    uint64
		expected medley.Node
	}{
		{
			ring: ring{
				{Node: "value1", Value: 100},
			},
			value:    50,
			expected: "value1",
		},
		{
			ring: ring{
				{Node: "value1", Value: 100},
			},
			value:    100,
			expected: "value1",
		},
		{
			ring: ring{
				{Node: "value1", Value: 100},
			},
			value:    200,
			expected: "value1",
		},
		{
			ring: ring{
				{Node: "value4", Value: 400}, {Node: "value3", Value: 300}, {Node: "value1", Value: 100}, {Node: "value5", Value: 500}, {Node: "value2", Value: 200}, {Node: "collision", Value: 200},
			},
			value:    50,
			expected: "value1",
		},
		{
			ring: ring{
				{Node: "value4", Value: 400}, {Node: "value3", Value: 300}, {Node: "value1", Value: 100}, {Node: "value5", Value: 500}, {Node: "value2", Value: 200}, {Node: "collision", Value: 200},
			},
			value:    150,
			expected: "collision",
		},
		{
			ring: ring{
				{Node: "value4", Value: 400}, {Node: "value3", Value: 300}, {Node: "value1", Value: 100}, {Node: "value5", Value: 500}, {Node: "value2", Value: 200}, {Node: "collision", Value: 200},
			},
			value:    250,
			expected: "value3",
		},
		{
			ring: ring{
				{Node: "value4", Value: 400}, {Node: "value3", Value: 300}, {Node: "value1", Value: 100}, {Node: "value5", Value: 500}, {Node: "value2", Value: 200}, {Node: "collision", Value: 200},
			},
			value:    1000,
			expected: "value1",
		},
	}

	for i, record := range testData {
		suite.Run(strconv.Itoa(i), func() {
			record.ring.sort()
			actual, err := record.ring.closest(record.value)
			suite.NoError(err)
			suite.Equal(record.expected, actual)
		})
	}
}

func TestRing(t *testing.T) {
	suite.Run(t, new(RingTestSuite))
}
