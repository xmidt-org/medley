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

func TestService(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}
