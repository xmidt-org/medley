// SPDX-FileCopyrightText: 2026 Comcast Cable Communications Management, LLC
// SPDX-License-Identifier: Apache-2.0

package consistent

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type RingSuite struct {
	suite.Suite
}

func TestRing(t *testing.T) {
	suite.Run(t, new(RingSuite))
}
