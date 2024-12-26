package consistent

import (
	"hash/fnv"
	"sort"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/xmidt-org/medley"
)

type BuilderSuite struct {
	suite.Suite

	object []byte
}

func (suite *BuilderSuite) SetupTest() {
	suite.object = []byte("test value")
}

func (suite *BuilderSuite) testStringsDefault() {
	services := []string{"service1", "service2", "service3"}
	ring := Strings(services...).Build()
	suite.Require().NotNil(ring)
	suite.Require().True(sort.IsSorted(ring.nodes))

	result, err := ring.Find(suite.object)
	suite.NoError(err)
	suite.Contains(services, result)
}

func (suite *BuilderSuite) testStringsCustom() {
	services := []string{"service1", "service2", "service3", "additional1", "additional2"}
	ring := Strings(services[:3]...).
		VNodes(100).
		Services(services[3:]...).
		Algorithm(medley.Algorithm{New64: fnv.New64}).
		ServiceHasher(nil). // force the default
		Build()

	suite.Require().NotNil(ring)
	suite.Require().True(sort.IsSorted(ring.nodes))

	result, err := ring.Find(suite.object)
	suite.NoError(err)
	suite.Contains(services, result)
}

func (suite *BuilderSuite) TestStrings() {
	suite.Run("Default", suite.testStringsDefault)
	suite.Run("Custom", suite.testStringsCustom)
}

func (suite *BuilderSuite) TestServices() {
	services := []string{"service1", "service2", "service3"}
	ring := Services(services...).Build()
	suite.Require().NotNil(ring)
	suite.Require().True(sort.IsSorted(ring.nodes))

	result, err := ring.Find(suite.object)
	suite.NoError(err)
	suite.Contains(services, result)
}

func (suite *BuilderSuite) TestBasicServices() {
	services := []medley.BasicService{
		{Host: "service1.net"},
		{Host: "service2.net", Port: 8080},
		{Host: "service3.net", Path: "/foo/bar"},
	}

	ring := BasicServices(services...).Build()
	suite.Require().NotNil(ring)
	suite.Require().True(sort.IsSorted(ring.nodes))

	result, err := ring.Find(suite.object)
	suite.NoError(err)
	suite.Contains(services, result)
}

func TestBuilder(t *testing.T) {
	suite.Run(t, new(BuilderSuite))
}
