package consistent

import (
	"hash/fnv"
	"strconv"
	"testing"

	"github.com/spaolacci/murmur3"
	"github.com/stretchr/testify/suite"
	"github.com/xmidt-org/medley"
)

type HashTestSuite struct {
	suite.Suite

	// goodConfigs are the configurations we execute
	// behavior tests with
	goodConfigs []Config
}

var _ suite.SetupTestSuite = (*HashTestSuite)(nil)

func (suite *HashTestSuite) SetupTest() {
	suite.goodConfigs = []Config{
		{},
		{
			Algorithm: medley.AlgorithmFNV,
			Vnodes:    10,
		},
		{
			Algorithm: "custom",
			Vnodes:    10,
			Extensions: map[string]medley.Algorithm{
				"custom": murmur3.New64,
			},
		},
	}
}

func (suite *HashTestSuite) TestNew() {
	suite.Run("Success", func() {
		testData := []struct {
			cfg            Config
			expectedVnodes int
		}{
			{
				cfg:            Config{},
				expectedVnodes: DefaultVnodes,
			},
			{
				cfg: Config{
					Algorithm: medley.AlgorithmFNV,
					Vnodes:    -123,
				},
				expectedVnodes: DefaultVnodes,
			},
			{
				cfg: Config{
					Algorithm: "custom",
					Vnodes:    42,
					Extensions: map[string]medley.Algorithm{
						"custom": fnv.New64a,
					},
				},
				expectedVnodes: 42,
			},
		}

		for i, record := range testData {
			suite.Run(strconv.Itoa(i), func() {
				h, err := New(record.cfg)
				suite.Require().NoError(err)
				suite.Require().NotNil(h)

				suite.NotNil(h.Algorithm())
				suite.Equal(record.expectedVnodes, h.Vnodes())
				suite.Zero(h.Len())
			})
		}
	})

	suite.Run("Fail", func() {
		testData := []Config{
			{
				Algorithm: "gblarglefargle",
			},
			{
				Algorithm: "stilldoesnotexist",
				Extensions: map[string]medley.Algorithm{
					"custom": fnv.New64a,
				},
			},
		}

		for i, cfg := range testData {
			suite.Run(strconv.Itoa(i), func() {
				h, err := New(cfg)
				suite.Error(err)
				suite.Nil(h)
			})
		}
	})
}

func (suite *HashTestSuite) TestGetAddRemove() {
	for i, cfg := range suite.goodConfigs {
		suite.Run(strconv.Itoa(i), func() {
			h, err := New(cfg)
			suite.Require().NoError(err)
			suite.Require().NotNil(h)

			suite.Equal(0, h.Len())
			_, err = h.Get(medley.String("key"))
			suite.Error(err)

			suite.T().Log("adding one hostname node")
			suite.Equal(1, h.Add([]medley.Node{"hostname-first.com"}))
			suite.Equal(1, h.Len())
			n, err := h.Get(medley.String("key"))
			suite.NoError(err)
			suite.Equal(medley.Node("hostname-first.com"), n)

			// idempotent
			suite.T().Log("adding the first hostname node again shouldn't change anything")
			suite.Equal(0, h.Add([]medley.Node{"hostname-first.com"}))
			suite.Equal(1, h.Len())
			n, err = h.Get(medley.String("key"))
			suite.NoError(err)
			suite.Equal(medley.Node("hostname-first.com"), n)

			suite.T().Log("adding a second hostname node")
			suite.Equal(1, h.Add([]medley.Node{"hostname-second.com"}))
			suite.Equal(2, h.Len())
			n, err = h.Get(medley.String("key"))
			suite.NoError(err)
			expected, err := h.ring.closest(medley.ComputeHash(medley.String("key"), h.alg))
			suite.Require().NoError(err)
			suite.Equal(expected, n)

			// idempotent
			suite.T().Log("adding the second hostname node again shouldn't change anything")
			suite.Equal(0, h.Add([]medley.Node{"hostname-second.com"}))
			suite.Equal(2, h.Len())
			n, err = h.Get(medley.String("key"))
			suite.NoError(err)
			expected, err = h.ring.closest(medley.ComputeHash(medley.String("key"), h.alg))
			suite.Require().NoError(err)
			suite.Equal(expected, n)

			suite.T().Log("removing the first hostname node")
			suite.Equal(1, h.Remove([]medley.Node{"hostname-first.com"}))
			suite.Equal(1, h.Len())
			n, err = h.Get(medley.String("key"))
			suite.NoError(err)
			suite.Equal(medley.Node("hostname-second.com"), n)

			// idempotent
			suite.T().Log("removing the first hostname node shouldn't change anything")
			suite.Equal(0, h.Remove([]medley.Node{"hostname-first.com"}))
			suite.Equal(1, h.Len())
			n, err = h.Get(medley.String("key"))
			suite.NoError(err)
			suite.Equal(medley.Node("hostname-second.com"), n)

			suite.T().Log("removing the second hostname node, which puts the hash back at empty")
			suite.Equal(1, h.Remove([]medley.Node{"hostname-second.com"}))
			suite.Equal(0, h.Len())
			_, err = h.Get(medley.String("key"))
			suite.Error(err)
		})
	}
}

func (suite *HashTestSuite) TestRehash() {
	for i, cfg := range suite.goodConfigs {
		suite.Run(strconv.Itoa(i), func() {
			h, err := New(cfg)
			suite.Require().NoError(err)
			suite.Require().NotNil(h)

			suite.T().Log("adding the initial set of nodes")
			suite.Require().Equal(
				2,
				h.Add([]medley.Node{"hostname-first.com", "hostname-second.com"}),
			)

			suite.Require().Equal(2, h.Len())

			n, err := h.Get(medley.String("key"))
			suite.NoError(err)
			expected, err := h.ring.closest(medley.ComputeHash(medley.String("key"), h.alg))
			suite.NoError(err)
			suite.Equal(expected, n)

			suite.T().Log("rehashing with the same nodes, which shouldn't change anything")
			added, removed := h.Rehash([]medley.Node{"hostname-second.com", "hostname-first.com"})
			suite.Equal(2, h.Len())
			suite.Zero(added)
			suite.Zero(removed)
			n, err = h.Get(medley.String("key"))
			suite.NoError(err)
			expected, err = h.ring.closest(medley.ComputeHash(medley.String("key"), h.alg))
			suite.NoError(err)
			suite.Equal(expected, n)
			suite.Contains(
				[]medley.Node{"hostname-second.com", "hostname-first.com"},
				n,
			)

			suite.T().Log("rehashing with an intersecting set")
			added, removed = h.Rehash(
				[]medley.Node{"hostname-third.com", "hostname-first.com", "hostname-fourth.com"},
			)
			suite.Equal(3, h.Len())
			suite.Equal(2, added)
			suite.Equal(1, removed)
			n, err = h.Get(medley.String("key"))
			suite.NoError(err)
			expected, err = h.ring.closest(medley.ComputeHash(medley.String("key"), h.alg))
			suite.NoError(err)
			suite.Equal(expected, n)
			suite.Contains(
				[]medley.Node{"hostname-third.com", "hostname-first.com", "hostname-fourth.com"},
				n,
			)

			suite.T().Log("rehashing with a disjoint set")
			added, removed = h.Rehash(
				[]medley.Node{"disjoint-1.net", "disjoint-2.net", "disjoint-3.net", "disjoint-4.net"},
			)
			suite.Equal(4, h.Len())
			suite.Equal(4, added)
			suite.Equal(3, removed)
			n, err = h.Get(medley.String("key"))
			suite.NoError(err)
			expected, err = h.ring.closest(medley.ComputeHash(medley.String("key"), h.alg))
			suite.NoError(err)
			suite.Equal(expected, n)
			suite.Contains(
				[]medley.Node{"disjoint-1.net", "disjoint-2.net", "disjoint-3.net", "disjoint-4.net"},
				n,
			)
		})
	}
}

func TestHash(t *testing.T) {
	suite.Run(t, new(HashTestSuite))
}
