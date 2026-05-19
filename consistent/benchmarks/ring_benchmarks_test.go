// SPDX-FileCopyrightText: 2026 Comcast Cable Communications Management, LLC
// SPDX-License-Identifier: Apache-2.0

package benchmarks

import (
	"fmt"
	"slices"
	"strconv"
	"testing"

	"github.com/billhathaway/consistentHash"
	"github.com/spaolacci/murmur3"
	"github.com/xmidt-org/medley"
	"github.com/xmidt-org/medley/consistent"
)

type benchmarkCase struct {
	name      string
	vnodes    int
	hostNames []string
}

var benchmarkCases []benchmarkCase

func init() {
	benchmarkHostNames := make([]string, 100)
	for i := range 100 {
		benchmarkHostNames[i] = fmt.Sprintf("benchmark-%d.medley.benchmarks.net", i)
	}

	benchmarkVnodes := []int{1, 50, 100, 200, 500}
	benchmakrHostCounts := []int{1, 10, 100}

	// multiply by 2 because for each combo, we do 2 algorithms:  default, and FNV64a
	benchmarkCases = make([]benchmarkCase, 0, len(benchmarkVnodes)*len(benchmakrHostCounts))

	for _, vnodes := range benchmarkVnodes {
		for _, hostCount := range benchmakrHostCounts {
			benchmarkCases = append(benchmarkCases, benchmarkCase{
				name:      fmt.Sprintf("vnodes=%03d|hostCount=%03d", vnodes, hostCount),
				vnodes:    vnodes,
				hostNames: benchmarkHostNames[:hostCount],
			})
		}
	}
}

func BenchmarkMultipleHashWrites(b *testing.B) {
	h := murmur3.New64()
	var iBuffer [10]byte
	for b.Loop() {
		h.Reset()
		h.Write(strconv.AppendUint(iBuffer[:], 204, 10))
		h.Write([]byte{'='})
		h.Write([]byte("hostname.something.net"))

		r := h.Sum64()
		r++ // so the compiler doesn't optimize it away
	}
}

func BenchmarkSingleHashWrite(b *testing.B) {
	h := murmur3.New64()
	const hostName = "hostname.something.net"
	for b.Loop() {
		h.Reset()
		h.Write(
			[]byte(strconv.Itoa(204) + "=" + hostName),
		)

		r := h.Sum64()
		r++ // so the compiler doesn't optimize it away
	}
}

// BenchmarkMedleyRingCreationUsingExpectedValues tests using a Builder to create
// a Ring with the ExpectedValues hint to preallocate the ring.
func BenchmarkMedleyRingCreationUsingExpectedValues(b *testing.B) {
	for _, benchmarkCase := range benchmarkCases {
		b.Run(benchmarkCase.name, func(b *testing.B) {
			builder := new(consistent.Builder[string]).
				VNodes(benchmarkCase.vnodes).
				ExpectedValues(len(benchmarkCase.hostNames))

			values := medley.Objectify(
				medley.String,
				slices.Values(benchmarkCase.hostNames),
			)

			for b.Loop() {
				builder.Build(values)
			}
		})
	}
}

// BenchmarkMedleyRingCreationNoExpectedValues tests using a Builder to create
// a Ring with no hint to preallocate.
func BenchmarkMedleyRingCreationNoExpectedValues(b *testing.B) {
	for _, benchmarkCase := range benchmarkCases {
		b.Run(benchmarkCase.name, func(b *testing.B) {
			builder := new(consistent.Builder[string]).
				VNodes(benchmarkCase.vnodes)

			values := medley.Objectify(
				medley.String,
				slices.Values(benchmarkCase.hostNames),
			)

			for b.Loop() {
				builder.Build(values)
			}
		})
	}
}

// BenchmarkConsistentHashCreation uses the consistenthash package to create
// the same rings as BenchmarkMedleyRingCreationUsingExpectedValues and
// BenchmarkMedleyRingCreationNoExpectedValues.
func BenchmarkConsistentHashCreation(b *testing.B) {
	for _, benchmarkCase := range benchmarkCases {
		b.Run(benchmarkCase.name, func(b *testing.B) {
			for b.Loop() {
				// apple-to-apples:
				// the Builder.Build does what this code does
				ch := consistentHash.New()
				ch.SetVnodeCount(benchmarkCase.vnodes)

				for _, hostName := range benchmarkCase.hostNames {
					ch.Add(hostName)
				}
			}
		})
	}
}
