package medley

import (
	"hash/fnv"
	"strconv"
	"testing"

	"github.com/spaolacci/murmur3"
	"github.com/stretchr/testify/suite"
)

type AlgorithmTestSuite struct {
	suite.Suite
}

func (suite *AlgorithmTestSuite) TestGetAlgorithm() {
	testData := []struct {
		name      string
		expectErr bool
	}{
		{
			name:      "",
			expectErr: false,
		},
		{
			name:      AlgorithmFNV,
			expectErr: false,
		},
		{
			name:      AlgorithmMurmur3,
			expectErr: false,
		},
		{
			name:      "unknown",
			expectErr: true,
		},
	}

	for i, record := range testData {
		suite.Run(strconv.Itoa(i), func() {
			alg, err := GetAlgorithm(record.name)
			suite.Equal(record.expectErr, alg == nil)
			if err != nil {
				suite.NotEmpty(err.Error())
			}
		})
	}
}

func (suite *AlgorithmTestSuite) TestFindAlgorithm() {
	testData := []struct {
		name       string
		extensions map[string]Algorithm
		expectErr  bool
	}{
		{
			name:       "",
			extensions: nil,
			expectErr:  false,
		},
		{
			name:       AlgorithmMurmur3,
			extensions: nil,
			expectErr:  false,
		},
		{
			name:       AlgorithmMurmur3,
			extensions: map[string]Algorithm{AlgorithmMurmur3: murmur3.New64},
			expectErr:  false,
		},
		{
			name:       AlgorithmMurmur3,
			extensions: map[string]Algorithm{"new": fnv.New64a},
			expectErr:  false,
		},
		{
			name:       "unknown",
			extensions: nil,
			expectErr:  true,
		},
		{
			name:       "unknown",
			extensions: map[string]Algorithm{"new": fnv.New64a},
			expectErr:  true,
		},
		{
			name:       "new",
			extensions: map[string]Algorithm{"new": fnv.New64a},
			expectErr:  false,
		},
	}

	for i, record := range testData {
		suite.Run(strconv.Itoa(i), func() {
			alg, err := FindAlgorithm(record.name, record.extensions)
			suite.Equal(record.expectErr, alg == nil)
			if err != nil {
				suite.NotEmpty(err.Error())
			}
		})
	}
}

func TestAlgorithm(t *testing.T) {
	suite.Run(t, new(AlgorithmTestSuite))
}
