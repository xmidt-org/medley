// SPDX-FileCopyrightText: 2025 Comcast Cable Communications Management, LLC
// SPDX-License-Identifier: Apache-2.0

package consistent

import (
	"fmt"
	"testing"

	"github.com/billhathaway/consistentHash"
)

var benchmarkVnodes = []int{50, 100, 200}

func BenchmarkRingCreation(b *testing.B) {
	for _, vnodes := range benchmarkVnodes {
		b.Run(
			fmt.Sprintf("vnodes-%d", vnodes),
			func(b *testing.B) {
				for range b.N {
					Strings(services[:]...).VNodes(vnodes).Build()
				}
			},
		)
	}
}

func BenchmarkConsistentHashCreation(b *testing.B) {
	for _, vnodes := range benchmarkVnodes {
		b.Run(
			fmt.Sprintf("vnodes-%d", vnodes),
			func(b *testing.B) {
				for range b.N {
					ch := consistentHash.New()
					ch.SetVnodeCount(vnodes)
					for _, service := range services {
						ch.Add(service)
					}
				}
			},
		)
	}
}
