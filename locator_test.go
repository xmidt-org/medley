// SPDX-FileCopyrightText: 2025 Comcast Cable Communications Management, LLC
// SPDX-License-Identifier: Apache-2.0

package medley

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type LocatorSuite struct {
	suite.Suite

	object       []byte
	objectString string
}

func (suite *LocatorSuite) find(ml *MultiLocator[string]) ([]string, error) {
	return ml.Find(suite.object)
}

func (suite *LocatorSuite) findString(ml *MultiLocator[string]) ([]string, error) {
	return ml.FindString(suite.objectString)
}

func (suite *LocatorSuite) SetupTest() {
	suite.objectString = "test value"
	suite.object = []byte(suite.objectString)
}

func (suite *LocatorSuite) assertExpectations(testObjects ...any) bool {
	return mock.AssertExpectationsForObjects(
		suite.T(),
		testObjects...,
	)
}

func (suite *LocatorSuite) TestFindString() {
	l := new(MockLocator[string])
	l.ExpectFindSuccess(suite.object, "service1").Once()

	actual, err := FindString(l, suite.objectString)
	suite.NoError(err)
	suite.Equal("service1", actual)

	suite.assertExpectations(l)
}

func (suite *LocatorSuite) testMultiLocatorFindEmpty() {
	ml := new(MultiLocator[string])
	results, err := ml.Find(suite.object)
	suite.ErrorIs(err, ErrNoServices)
	suite.Empty(results)
}

func (suite *LocatorSuite) testMultiLocatorFindStringEmpty() {
	ml := new(MultiLocator[string])
	results, err := ml.FindString(suite.objectString)
	suite.ErrorIs(err, ErrNoServices)
	suite.Empty(results)
}

// testMultiLocatorAllSuccess sets the expectation for a suite.object call on contained
// locators and lets a test pass in a closure to invoke a method on the MultiLocator under test.
func (suite *LocatorSuite) testMultiLocatorAllSuccess(finder func(*MultiLocator[string]) ([]string, error)) func() {
	return func() {
		var (
			l1 = new(MockLocator[string])
			l2 = new(MockLocator[string])
			l3 = new(MockLocator[string])

			ml = NewMultiLocator(l1, l2)
		)

		l1.ExpectFindSuccess(suite.object, "service1").Times(2)
		l2.ExpectFindSuccess(suite.object, "service2").Times(3)
		l3.ExpectFindSuccess(suite.object, "service3").Times(2)

		results, err := finder(ml)
		suite.NoError(err)
		suite.ElementsMatch([]string{"service1", "service2"}, results)

		ml.Add(l3)
		results, err = finder(ml)
		suite.NoError(err)
		suite.ElementsMatch([]string{"service1", "service2", "service3"}, results)

		ml.Remove(l1)
		results, err = finder(ml)
		suite.NoError(err)
		suite.ElementsMatch([]string{"service2", "service3"}, results)

		suite.assertExpectations(l1, l2, l3)
	}
}

// testMultiLocatorSomeMissingServices tests that FindXXX works correctly when some locators
// are missing services, but others aren't.
func (suite *LocatorSuite) testMultiLocatorSomeMissingServices(finder func(*MultiLocator[string]) ([]string, error)) func() {
	return func() {
		var (
			l1 = new(MockLocator[string])
			l2 = new(MockLocator[string])
			l3 = new(MockLocator[string])

			ml = NewMultiLocator(l1, l2, l3)
		)

		l1.ExpectFindSuccess(suite.object, "service1").Once()
		l2.ExpectFindNoServices(suite.object).Once()
		l3.ExpectFindSuccess(suite.object, "service3").Once()

		results, err := finder(ml)
		suite.NoError(err)
		suite.ElementsMatch([]string{"service1", "service3"}, results)

		suite.assertExpectations(l1, l2, l3)
	}
}

// testMultiLocatorFail tests that FindXXX works correctly when a locator returns an error.
func (suite *LocatorSuite) testMultiLocatorFail(finder func(*MultiLocator[string]) ([]string, error)) func() {
	return func() {
		var (
			expectedErr = errors.New("expected")

			l1 = new(MockLocator[string])
			l2 = new(MockLocator[string])
			l3 = new(MockLocator[string])

			ml = NewMultiLocator(l1, l2, l3)
		)

		l1.ExpectFindSuccess(suite.object, "service1").Once()
		l2.ExpectFindFail(suite.object, expectedErr).Once()

		results, err := finder(ml)
		suite.ErrorIs(err, expectedErr)
		suite.Empty(results)

		suite.assertExpectations(l1, l2, l3)
	}
}

func (suite *LocatorSuite) TestMultiLocator() {
	suite.Run("Find", func() {
		suite.Run("Empty", suite.testMultiLocatorFindEmpty)
		suite.Run("AllSuccess", suite.testMultiLocatorAllSuccess(suite.find))
		suite.Run("SomeMissingServices", suite.testMultiLocatorSomeMissingServices(suite.find))
		suite.Run("Fail", suite.testMultiLocatorFail(suite.find))
	})

	suite.Run("FindString", func() {
		suite.Run("Empty", suite.testMultiLocatorFindStringEmpty)
		suite.Run("AllSuccess", suite.testMultiLocatorAllSuccess(suite.findString))
		suite.Run("SomeMissingServices", suite.testMultiLocatorSomeMissingServices(suite.findString))
		suite.Run("Fail", suite.testMultiLocatorFail(suite.findString))
	})
}

func (suite *LocatorSuite) TestUpdatableLocator() {
	var (
		expectedErr = errors.New("expected error")

		l1 = new(MockLocator[string])
		l2 = new(MockLocator[string])
		l3 = new(MockLocator[string])

		ul = NewUpdatableLocator(l1)
	)

	suite.Require().NotNil(ul)

	l1.ExpectFindSuccess(suite.object, "service1").Once()
	l2.ExpectFindSuccess(suite.object, "service2").Once()
	l3.ExpectFindFail(suite.object, expectedErr).Once()

	result, err := ul.Find(suite.object)
	suite.NoError(err)
	suite.Equal("service1", result)

	ul.Set(nil)
	result, err = ul.Find(suite.object)
	suite.ErrorIs(err, ErrNoServices)
	suite.Empty(result)

	ul.Set(l2)
	result, err = ul.Find(suite.object)
	suite.NoError(err)
	suite.Equal("service2", result)

	ul.Set(l3)
	result, err = ul.Find(suite.object)
	suite.ErrorIs(err, expectedErr)
	suite.Empty(result)

	suite.assertExpectations(l1, l2, l3)
}

func (suite *LocatorSuite) TestSetLocator() {
	var (
		l1 = new(MockLocator[string])
		l2 = new(MockLocator[string])
		ul = NewUpdatableLocator[string](nil)
	)

	l2.ExpectFindSuccess(suite.object, "service2").Once()
	suite.Require().NotNil(ul)

	suite.NotPanics(func() {
		SetLocator(l1, l2) // nop, since l1 doesn't have a Set method
	})

	SetLocator(ul, l2)
	result, err := ul.Find(suite.object)
	suite.NoError(err)
	suite.Equal("service2", result)

	suite.assertExpectations(l1, l2)
}

func TestLocator(t *testing.T) {
	suite.Run(t, new(LocatorSuite))
}
