package medley

import (
	"bytes"
	"strconv"
	"testing"

	"github.com/stretchr/testify/suite"
)

type NodeTestSuite struct {
	suite.Suite
}

func (suite *NodeTestSuite) TestWriteTo() {
	var output bytes.Buffer
	c, err := Node("test").WriteTo(&output)
	suite.Equal(int64(len("test")), c)
	suite.NoError(err)
}

func TestNode(t *testing.T) {
	suite.Run(t, new(NodeTestSuite))
}

type NodeSetTestSuite struct {
	suite.Suite
}

func (suite *NodeSetTestSuite) TestNewNodeSet() {
	testData := [][]Node{
		nil,
		{},
		{"value1"},
		{"value1", "value2"},
		{"value1", "value2", "value3", "value4", "value5"},
	}

	for i, nodes := range testData {
		suite.Run(strconv.Itoa(i), func() {
			nodeSet := NewNodeSet(nodes...)
			suite.Equal(len(nodes), len(nodeSet))
			suite.Equal(len(nodes), nodeSet.Len())

			for _, n := range nodes {
				suite.True(nodeSet.Has(n))
			}

			suite.False(nodeSet.Has("nosuch"))
		})
	}
}

func (suite *NodeSetTestSuite) TestAdd() {
	testData := []NodeSet{
		nil,
		{},
		NewNodeSet("value1"),
		NewNodeSet("value1", "value2"),
		NewNodeSet("value1", "value2", "value3", "value4", "value5"),
	}

	for i, nodeSet := range testData {
		suite.Run(strconv.Itoa(i), func() {
			existing := make([]Node, 0, nodeSet.Len())
			for k := range nodeSet {
				existing = append(existing, k)
			}

			suite.False(nodeSet.Has("test"))
			suite.True(nodeSet.Add("test"))
			suite.True(nodeSet.Has("test"))
			suite.Equal(nodeSet.Len(), len(existing)+1)
			for _, k := range existing {
				suite.True(nodeSet.Has(k))
			}

			// idempotency:
			suite.False(nodeSet.Add("test"))
			suite.True(nodeSet.Has("test"))
			suite.Equal(nodeSet.Len(), len(existing)+1)
			for _, k := range existing {
				suite.True(nodeSet.Has(k))
			}
		})
	}
}

func (suite *NodeSetTestSuite) TestAddAll() {
	testData := []struct {
		nodeSet       NodeSet
		nodes         []Node
		expected      NodeSet
		expectedCount int
	}{
		{
			nodeSet:       nil,
			nodes:         nil,
			expected:      nil,
			expectedCount: 0,
		},
		{
			nodeSet:       NodeSet{},
			nodes:         []Node{},
			expected:      NodeSet{},
			expectedCount: 0,
		},
		{
			nodeSet:       nil,
			nodes:         []Node{"value1"},
			expected:      NewNodeSet("value1"),
			expectedCount: 1,
		},
		{
			nodeSet:       NodeSet{},
			nodes:         []Node{"value1"},
			expected:      NewNodeSet("value1"),
			expectedCount: 1,
		},
		{
			nodeSet:       NewNodeSet("value1"),
			nodes:         []Node{"value1"},
			expected:      NewNodeSet("value1"),
			expectedCount: 0,
		},
		{
			nodeSet:       NewNodeSet("value1", "value2"),
			nodes:         []Node{"value1", "value2"},
			expected:      NewNodeSet("value1", "value2"),
			expectedCount: 0,
		},
		{
			nodeSet:       NewNodeSet("value1", "value2"),
			nodes:         []Node{"value1"},
			expected:      NewNodeSet("value1", "value2"),
			expectedCount: 0,
		},
		{
			nodeSet:       NewNodeSet("value1", "value3", "value5"),
			nodes:         []Node{"value2", "value4", "value6"},
			expected:      NewNodeSet("value1", "value2", "value3", "value4", "value5", "value6"),
			expectedCount: 3,
		},
	}

	for i, record := range testData {
		suite.Run(strconv.Itoa(i), func() {
			initialLen := record.nodeSet.Len()
			suite.Equal(
				record.expectedCount,
				record.nodeSet.AddAll(record.nodes...),
			)

			suite.Equal(
				initialLen+record.expectedCount,
				record.nodeSet.Len(),
			)

			for _, n := range record.nodes {
				suite.True(record.nodeSet.Has(n))
			}

			// idempotency:
			suite.Zero(record.nodeSet.AddAll(record.nodes...))

			suite.Equal(
				initialLen+record.expectedCount,
				record.nodeSet.Len(),
			)

			for _, n := range record.nodes {
				suite.True(record.nodeSet.Has(n))
			}
		})
	}
}

func (suite *NodeSetTestSuite) TestRemove() {
	testData := []NodeSet{
		nil,
		{},
		NewNodeSet("value1"),
		NewNodeSet("value1", "test"),
		NewNodeSet("value1", "value2"),
		NewNodeSet("value1", "value2", "test"),
		NewNodeSet("value1", "value2", "value3", "value4", "value5"),
		NewNodeSet("value1", "value2", "value3", "value4", "value5", "test"),
	}

	for i, nodeSet := range testData {
		suite.Run(strconv.Itoa(i), func() {
			var (
				initialLen  = nodeSet.Len()
				expectedLen int
				exists      = nodeSet.Has("test")
			)

			if exists {
				expectedLen = initialLen - 1
			} else {
				expectedLen = initialLen
			}

			suite.Equal(exists, nodeSet.Remove("test"))
			suite.Equal(expectedLen, nodeSet.Len())
			suite.False(nodeSet.Has("test"))

			// idempotency:
			suite.False(nodeSet.Remove("test"))
			suite.Equal(expectedLen, nodeSet.Len())
			suite.False(nodeSet.Has("test"))
		})
	}
}

func (suite *NodeSetTestSuite) TestRemoveAll() {
	testData := []struct {
		nodeSet       NodeSet
		nodes         []Node
		expected      NodeSet
		expectedCount int
	}{
		{
			nodeSet:       nil,
			nodes:         nil,
			expected:      nil,
			expectedCount: 0,
		},
		{
			nodeSet:       NodeSet{},
			nodes:         []Node{},
			expected:      NodeSet{},
			expectedCount: 0,
		},
		{
			nodeSet:       nil,
			nodes:         []Node{"value1"},
			expected:      NodeSet{},
			expectedCount: 0,
		},
		{
			nodeSet:       NodeSet{},
			nodes:         []Node{"value1"},
			expected:      NodeSet{},
			expectedCount: 0,
		},
		{
			nodeSet:       NewNodeSet("value1"),
			nodes:         []Node{"value1"},
			expected:      NodeSet{},
			expectedCount: 1,
		},
		{
			nodeSet:       NewNodeSet("value1", "value2", "value3", "value4", "value5"),
			nodes:         []Node{"value1"},
			expected:      NewNodeSet("value2", "value3", "value4", "value5"),
			expectedCount: 1,
		},
		{
			nodeSet:       NewNodeSet("value1", "value2", "value3", "value4", "value5"),
			nodes:         []Node{"value2", "value5"},
			expected:      NewNodeSet("value1", "value3", "value4"),
			expectedCount: 2,
		},
		{
			nodeSet:       NewNodeSet("value1", "value2", "value3", "value4", "value5"),
			nodes:         []Node{"value1", "value2", "value5"},
			expected:      NewNodeSet("value3", "value4"),
			expectedCount: 3,
		},
		{
			nodeSet:       NewNodeSet("value1", "value2", "value3", "value4", "value5"),
			nodes:         []Node{"value1", "value2", "value3", "value4", "value5"},
			expected:      NodeSet{},
			expectedCount: 5,
		},
	}

	for i, record := range testData {
		suite.Run(strconv.Itoa(i), func() {
			initialLen := record.nodeSet.Len()
			suite.Equal(
				record.expectedCount,
				record.nodeSet.RemoveAll(record.nodes...),
			)

			suite.Equal(
				initialLen-record.expectedCount,
				record.nodeSet.Len(),
			)

			for _, n := range record.nodes {
				suite.False(record.nodeSet.Has(n))
			}

			// idempotency:
			suite.Zero(record.nodeSet.RemoveAll(record.nodes...))

			suite.Equal(
				initialLen-record.expectedCount,
				record.nodeSet.Len(),
			)

			for _, n := range record.nodes {
				suite.False(record.nodeSet.Has(n))
			}
		})
	}
}

func (suite *NodeSetTestSuite) TestFilter() {
	testData := []struct {
		nodeSet NodeSet
		input   []Node
	}{
		{
			nodeSet: NodeSet{},
			input:   []Node{},
		},
		{
			nodeSet: NewNodeSet("server-1"),
			input:   []Node{"server-1"},
		},
		{
			nodeSet: NewNodeSet("server-1"),
			input:   []Node{"server-2", "server-1", "server-5"},
		},
		{
			nodeSet: NewNodeSet("server-1"),
			input:   []Node{"server-2", "server-6", "server-1", "server-5"},
		},
		{
			nodeSet: NewNodeSet("server-5", "server-267", "server-4", "server-12"),
			input:   []Node{},
		},
		{
			nodeSet: NewNodeSet("server-5", "server-267", "server-4", "server-12"),
			input:   []Node{"server-5"},
		},
		{
			nodeSet: NewNodeSet("server-5", "server-267", "server-4", "server-12"),
			input:   []Node{"server-16", "server-267", "server-22", "server-5", "server-12"},
		},
		{
			nodeSet: NewNodeSet("server-5", "server-267", "server-4", "server-12"),
			input:   []Node{"server-16", "server-267", "server-22", "server-5", "server-12", "server-44"},
		},
		{
			nodeSet: NewNodeSet("server-5", "server-267", "server-4", "server-12"),
			input:   []Node{"server-5", "server-267", "server-4", "server-12"},
		},
		{
			nodeSet: NewNodeSet("server-5", "server-267", "server-4", "server-12"),
			input:   []Node{"server-924", "server-34786283", "server-376", "server-234355"},
		},
	}

	for i, record := range testData {
		suite.Run(strconv.Itoa(i), func() {
			in, notIn := record.nodeSet.Filter(record.input)
			suite.Require().Equal(len(record.input), len(in)+len(notIn))

			suite.Equal(record.input[0:len(in)], in)
			suite.Equal(record.input[len(in):], notIn)

			for _, n := range in {
				suite.True(record.nodeSet[n])
			}

			for _, n := range notIn {
				suite.False(record.nodeSet[n])
			}
		})
	}
}

func TestNodeSet(t *testing.T) {
	suite.Run(t, new(NodeSetTestSuite))
}
