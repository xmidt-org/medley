package consistent

import (
	"testing"

	"github.com/spaolacci/murmur3"
	"github.com/stretchr/testify/suite"
	"github.com/xmidt-org/medley"
)

type AssignerTestSuite struct {
	suite.Suite
}

func (suite *AssignerTestSuite) TestResetAndNext() {
	a := newAssigner(murmur3.New64)
	suite.Require().NotNil(a)

	for _, node := range []medley.Node{"test1", "test2"} {
		a.reset(node)
		values := make(map[uint64]bool)
		for i := 0; i < 10; i++ {
			values[a.next()] = true
		}

		suite.Len(values, 10)
	}
}

func TestAssigner(t *testing.T) {
	suite.Run(t, new(AssignerTestSuite))
}
