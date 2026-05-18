// SPDX-FileCopyrightText: 2026 Comcast Cable Communications Management, LLC
// SPDX-License-Identifier: Apache-2.0

package medley

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type AlgorithmSuite[R uint32 | uint64] struct {
	suite.Suite
}

type Algorithm32Suite struct {
	AlgorithmSuite[uint32]
}

type Algorithm64Suite struct {
	AlgorithmSuite[uint64]
}

func TestAlgorithm32(t *testing.T) {
	suite.Run(t, new(Algorithm32Suite))
}

func TestAlgorithm64(t *testing.T) {
	suite.Run(t, new(Algorithm64Suite))
}
