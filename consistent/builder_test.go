// SPDX-FileCopyrightText: 2026 Comcast Cable Communications Management, LLC
// SPDX-License-Identifier: Apache-2.0

package consistent

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type BuilderTestSuite struct {
	suite.Suite
}

func TestBuilder(t *testing.T) {
	suite.Run(t, new(BuilderTestSuite))
}
