package medley

import (
	"bytes"
	"hash/fnv"
	"testing"

	"github.com/spaolacci/murmur3"
	"github.com/stretchr/testify/suite"
)

type KeyTestSuite struct {
	suite.Suite
}

func (suite *KeyTestSuite) TestComputeHash() {
	suite.NotZero(
		ComputeHash(String("key"), fnv.New64a),
	)

	suite.NotZero(
		ComputeHash(Bytes([]byte{1, 2, 3, 4, 5}), murmur3.New64),
	)
}

func (suite *KeyTestSuite) TestBytes() {
	key := Bytes([]byte{1, 2, 3, 4})
	var o bytes.Buffer

	c, err := key.WriteTo(&o)
	suite.NoError(err)
	suite.Equal(int64(len(key)), c)
	suite.Equal([]byte{1, 2, 3, 4}, o.Bytes())
}

func (suite *KeyTestSuite) TestString() {
	key := String("hello, world")
	var o bytes.Buffer

	c, err := key.WriteTo(&o)
	suite.NoError(err)
	suite.Equal(int64(len(key)), c)
	suite.Equal([]byte(key), o.Bytes())
}

func TestKey(t *testing.T) {
	suite.Run(t, new(KeyTestSuite))
}
