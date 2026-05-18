// SPDX-FileCopyrightText: 2026 Comcast Cable Communications Management, LLC
// SPDX-License-Identifier: Apache-2.0

package consistent

import (
	"fmt"
	"iter"
	"reflect"
	"slices"
	"sort"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/xmidt-org/medley"
)

type testService struct {
	hostName string
	port     int
}

type BuilderTestSuite struct {
	suite.Suite

	// testServices are some slices of input values for hashing
	testServices [][]testService
}

// generateTestServices creates some testService instances that are different enough to
// exercize the hashing logic.
func (suite *BuilderTestSuite) generateTestServices() []testService {
	return []testService{
		{hostName: "test-123.test.org", port: 1010},
		{hostName: "manic.thurston.net", port: 9180},
		{hostName: "something.fast-api.somethingelse.net", port: 70},
		{hostName: "gargantuan.medley.net", port: 746},
		{hostName: "randomizer.host.com", port: 1290},
		{hostName: "aggregate.yahoo.com", port: 6504},
		{hostName: "amazingly.fast.org", port: 1400},
	}
}

func (suite *BuilderTestSuite) SetupSuite() {
	allTestServices := suite.generateTestServices()
	suite.testServices = append(suite.testServices, []testService{})
	suite.testServices = append(suite.testServices, allTestServices[:2])
	suite.testServices = append(suite.testServices, allTestServices[:5])
	suite.testServices = append(suite.testServices, allTestServices)
}

// values builds the sequence that Builder.Build expects from a slice of test services.
// The hostName is used as the hashing object.
func (suite *BuilderTestSuite) values(services []testService) iter.Seq2[medley.Object, testService] {
	return medley.Objectify(
		func(ts testService) medley.Object { return medley.String(ts.hostName) },
		slices.Values(services),
	)
}

// runBuildTests does the grunt work of executing 1 test per chunk of test services, using the
// given closure to create the Ring.
func (suite *BuilderTestSuite) runBuildTests(expectedVNodes int, ringer func([]testService) *Ring[testService]) {
	for _, testServices := range suite.testServices {
		suite.Run(fmt.Sprintf("values=%d", len(testServices)), func() {
			ring := ringer(testServices)
			suite.Require().NotNil(ring)
			suite.Require().NotNil(ring.alg) // the builder should always set this, even if it's the default
			suite.Require().Len(ring.nodes, expectedVNodes*len(testServices))
			suite.Require().True(sort.IsSorted(ring.nodes))
			if len(ring.nodes) == 0 {
				return
			}

			// take a set of test client names to hash to, and make sure they agree
			// with the hash nodes
			for _, clientName := range []string{"aclient", "homersimpson", "123anywhere"} {
				nearest := ring.Nearest(medley.String(clientName))
				suite.False(reflect.ValueOf(nearest).IsZero())
			}

			ring.Clear()
			suite.Len(ring.nodes, 0)
			ts := ring.Nearest(medley.String("aclient"))
			suite.True(reflect.ValueOf(ts).IsZero())
		})
	}
}

func (suite *BuilderTestSuite) testBuildDefault(testServices []testService) *Ring[testService] {
	var builder Builder[testService]
	return builder.Build(
		suite.values(testServices),
	)
}

func (suite *BuilderTestSuite) testBuildCustom(vnodes int, alg *medley.Algorithm[uint64]) func([]testService) *Ring[testService] {
	return func(testServices []testService) *Ring[testService] {
		var builder Builder[testService]
		builder.VNodes(vnodes).Algorithm(alg).ExpectedValues(len(testServices))
		return builder.Build(
			suite.values(testServices),
		)
	}
}

func (suite *BuilderTestSuite) TestBuild() {
	suite.Run("Default", func() {
		suite.runBuildTests(DefaultVNodes, suite.testBuildDefault)
	})

	suite.Run("CustomVNodes", func() {
		for _, vnodes := range []int{1, 10, 500} {
			suite.Run(fmt.Sprintf("vnodes=%d", vnodes), func() {
				suite.runBuildTests(vnodes, suite.testBuildCustom(vnodes, nil))
			})
		}
	})

	suite.Run("CustomAlgorithm", func() {
		suite.Run("FNV64", func() {
			suite.runBuildTests(DefaultVNodes, suite.testBuildCustom(DefaultVNodes, medley.FNV64()))
		})

		suite.Run("FNV64a", func() {
			suite.runBuildTests(DefaultVNodes, suite.testBuildCustom(DefaultVNodes, medley.FNV64a()))
		})
	})
}

func TestBuilder(t *testing.T) {
	suite.Run(t, new(BuilderTestSuite))
}
