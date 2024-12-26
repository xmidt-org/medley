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

func TestService(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}
