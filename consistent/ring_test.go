// SPDX-FileCopyrightText: 2025 Comcast Cable Communications Management, LLC
// SPDX-License-Identifier: Apache-2.0

package consistent

import (
	"sort"
	"testing"

	"github.com/billhathaway/consistentHash"
	"github.com/stretchr/testify/suite"
	"github.com/xmidt-org/medley"
)

type RingSuite struct {
	suite.Suite

	originalServices []string
	original         *Ring[string]
}

func (suite *RingSuite) SetupSuite() {
	// for the unit tests, we don't need all the services
	suite.originalServices = services[0:4]

	suite.original = Strings(suite.originalServices...).Build()
	suite.Require().NotNil(suite.original)
	suite.Require().True(sort.IsSorted(suite.original.nodes))

	distribution := make(map[string]int)
	for _, object := range hashObjects {
		result, err := suite.original.Find(object[:])
		suite.Require().NoError(err)
		suite.Require().Contains(suite.originalServices, result)
		distribution[result] += 1
	}

	// the distribution should be close to even
	expectedCount := objectCount / len(suite.originalServices)
	for _, actualCount := range distribution {
		// each count should be within 25% of its expected value.
		// 25% is just a guess, but it should prevent drift as
		// the codebase changes.
		suite.InEpsilon(expectedCount, actualCount, 0.25)
	}
}

func (suite *RingSuite) update(services ...string) (*Ring[string], bool) {
	updated, didUpdate := Update(suite.original, services...)
	suite.Require().NotNil(updated)
	suite.Require().True(sort.IsSorted(updated.nodes))

	return updated, didUpdate
}

func (suite *RingSuite) testUpdateEmpty() {
	// the list of updated services is empty
	updated, didUpdate := suite.update()
	suite.True(didUpdate)
	suite.Empty(updated.nodes)

	for _, object := range hashObjects {
		result, err := updated.Find(object[:])
		suite.Empty(result)
		suite.ErrorIs(err, medley.ErrNoServices)
	}
}

func (suite *RingSuite) testUpdatePartial() {
	partial := []string{"new1", suite.originalServices[0], "new2"}
	updated, didUpdate := suite.update(partial...)
	suite.True(didUpdate)

	for _, object := range hashObjects {
		result, err := updated.Find(object[:])
		suite.Contains(partial, result)
		suite.NoError(err)
	}
}

func (suite *RingSuite) testUpdateAllNew() {
	allNew := []string{"new1.service.com", "new2.service.com", "new3.service.com", "new4.service.com"}
	updated, didUpdate := suite.update(allNew...)
	suite.True(didUpdate)

	for _, object := range hashObjects {
		result, err := updated.Find(object[:])
		suite.Contains(allNew, result)
		suite.NoError(err)
	}
}

func (suite *RingSuite) testUpdateNotNeeded() {
	updated, didUpdate := suite.update(suite.originalServices...)
	suite.Same(suite.original, updated)
	suite.False(didUpdate)
}

func (suite *RingSuite) TestUpdate() {
	suite.Run("Empty", suite.testUpdateEmpty)
	suite.Run("Partial", suite.testUpdatePartial)
	suite.Run("AllNew", suite.testUpdateAllNew)
	suite.Run("NotNeeded", suite.testUpdateNotNeeded)
}

func (suite *RingSuite) TestBackwardCompatibility() {
	ch := consistentHash.New()
	ch.SetVnodeCount(DefaultVNodes)
	for _, service := range suite.originalServices {
		ch.Add(service)
	}

	var (
		oldResult string
		newResult string
		err       error
	)

	for _, object := range hashObjects {
		oldResult, err = ch.Get(object[:])
		suite.Require().NoError(err)

		newResult, err = suite.original.Find(object[:])
		suite.Require().NoError(err)

		suite.Equal(oldResult, newResult)
	}
}

func TestRing(t *testing.T) {
	suite.Run(t, new(RingSuite))
}
