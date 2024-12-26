package medley

import (
	"hash/fnv"
	"testing"

	"github.com/spaolacci/murmur3"
	"github.com/stretchr/testify/suite"
)

type AlgorithmSuite struct {
	suite.Suite

	hashInput string
	expected  uint64
}

// sum64 is just a Sum64 function that runs the fnv hash
func (suite *AlgorithmSuite) sum64(v []byte) uint64 {
	hash := fnv.New64()
	hash.Write(v)
	return hash.Sum64()
}

func (suite *AlgorithmSuite) SetupTest() {
	suite.hashInput = "test hash value"
	suite.expected = suite.sum64([]byte(suite.hashInput))
}

func (suite *AlgorithmSuite) assertExpected(v uint64) {
	suite.Equal(suite.expected, v)
}

func (suite *AlgorithmSuite) testSum64BytesUsingSum64() {
	alg := Algorithm{
		New64: fnv.New64,
		Sum64: suite.sum64,
	}

	suite.assertExpected(
		alg.Sum64Bytes([]byte(suite.hashInput)),
	)
}

func (suite *AlgorithmSuite) testSum64BytesUsingNew64() {
	alg := Algorithm{
		New64: fnv.New64,
		Sum64: nil,
	}

	suite.assertExpected(
		alg.Sum64Bytes([]byte(suite.hashInput)),
	)
}

func (suite *AlgorithmSuite) TestSum64Bytes() {
	suite.Run("UsingSum64", suite.testSum64BytesUsingSum64)
	suite.Run("UsingNew64", suite.testSum64BytesUsingNew64)
}

func (suite *AlgorithmSuite) testSum64StringUsingSum64() {
	alg := Algorithm{
		New64: fnv.New64,
		Sum64: suite.sum64,
	}

	suite.assertExpected(
		alg.Sum64String(suite.hashInput),
	)
}

func (suite *AlgorithmSuite) testSum64StringUsingNew64() {
	alg := Algorithm{
		New64: fnv.New64,
		Sum64: nil,
	}

	suite.assertExpected(
		alg.Sum64String(suite.hashInput),
	)
}

func (suite *AlgorithmSuite) TestSum64String() {
	suite.Run("UsingSum64", suite.testSum64StringUsingSum64)
	suite.Run("UsingNew64", suite.testSum64StringUsingNew64)
}

func (suite *AlgorithmSuite) TestDefaultAlgorithm() {
	alg := DefaultAlgorithm()
	suite.Require().NotNil(alg.New64)
	suite.Require().NotNil(alg.Sum64)

	expected := murmur3.Sum64([]byte(suite.hashInput))
	suite.Equal(expected, alg.Sum64Bytes([]byte(suite.hashInput)))
	suite.Equal(expected, alg.Sum64String(suite.hashInput))
}

func TestAlgorithm(t *testing.T) {
	suite.Run(t, new(AlgorithmSuite))
}
