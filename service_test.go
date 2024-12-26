// SPDX-FileCopyrightText: 2025 Comcast Cable Communications Management, LLC
// SPDX-License-Identifier: Apache-2.0

package medley

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/suite"
)

type ServiceSuite struct {
	suite.Suite
}

func (suite *ServiceSuite) TestDefaultServiceHasher() {
	var b bytes.Buffer
	suite.NoError(
		DefaultServiceHasher[string](&b, "service.com"),
	)

	suite.Equal("service.com", b.String())
}

func (suite *ServiceSuite) TestHashStringTo() {
	var b bytes.Buffer
	suite.NoError(
		HashStringTo(&b, "test value"),
	)

	suite.NotZero(b.Len())
}

func (suite *ServiceSuite) TestHashBasicServiceTo() {
	var b bytes.Buffer
	suite.NoError(
		HashBasicServiceTo(
			&b,
			BasicService{
				Scheme: "http",
				Host:   "service.com",
				Port:   1234,
				Path:   "/foo/bar",
			},
		),
	)

	suite.NotZero(b.Len())
}

func (suite *ServiceSuite) testMapUpdateNil() {
	var m Map[string, int] // nil Map
	suite.Zero(m.Len())

	expected := []string{"service1", "service2"}
	visited := make([]string, 0, len(expected))
	for u := range m.Update(expected...) {
		suite.False(u.Exists)
		suite.Zero(u.Value)
		visited = append(visited, u.Service)
	}

	suite.Equal(expected, visited)
}

func (suite *ServiceSuite) testMapUpdateAllServicesExist() {
	m := Map[string, int]{
		"service1": 123,
		"service2": 456,
		"service3": 789,
	}

	suite.Equal(3, m.Len())

	expected := []string{"service2", "service3"}
	visited := make([]string, 0, len(expected))
	for u := range m.Update(expected...) {
		suite.True(u.Exists)
		suite.Equal(m[u.Service], u.Value)
		visited = append(visited, u.Service)
	}

	suite.Equal(expected, visited)
}

func (suite *ServiceSuite) testMapUpdateNoServicesExist() {
	m := Map[string, int]{
		"service1": 123,
		"service2": 456,
		"service3": 789,
	}

	suite.Equal(3, m.Len())

	expected := []string{"service9", "service6", "service8"}
	visited := make([]string, 0, len(expected))
	for u := range m.Update(expected...) {
		suite.False(u.Exists)
		suite.Zero(u.Value)
		visited = append(visited, u.Service)
	}

	suite.Equal(expected, visited)
}

func (suite *ServiceSuite) testMapUpdateSomeServicesExist() {
	m := Map[string, int]{
		"service1": 123,
		"service2": 456,
		"service3": 789,
	}

	suite.Equal(3, m.Len())

	exists := make([]string, 0, 4)
	notExists := make([]string, 0, 4)
	for u := range m.Update("service4", "service3", "service8", "service1") {
		if u.Exists {
			suite.Equal(m[u.Service], u.Value)
			exists = append(exists, u.Service)
		} else {
			suite.Zero(u.Value)
			notExists = append(notExists, u.Service)
		}
	}

	suite.Equal([]string{"service3", "service1"}, exists)
	suite.Equal([]string{"service4", "service8"}, notExists)
}

func (suite *ServiceSuite) testMapUpdateBreak() {
	m := Map[string, int]{
		"service1": 123,
		"service2": 456,
		"service3": 789,
	}

	suite.Equal(3, m.Len())

	visited := make([]string, 0, 1)
	for u := range m.Update("first", "second") {
		if u.Service != "first" {
			break
		}

		visited = append(visited, u.Service)
	}

	suite.Equal([]string{"first"}, visited)
}

func (suite *ServiceSuite) TestMap() {
	suite.Run("Update", func() {
		suite.Run("Nil", suite.testMapUpdateNil)
		suite.Run("AllServicesExist", suite.testMapUpdateAllServicesExist)
		suite.Run("NoServicesExist", suite.testMapUpdateNoServicesExist)
		suite.Run("SomeServicesExist", suite.testMapUpdateSomeServicesExist)
		suite.Run("Break", suite.testMapUpdateBreak)
	})
}

func TestService(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}
