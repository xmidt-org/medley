// SPDX-FileCopyrightText: 2025 Comcast Cable Communications Management, LLC
// SPDX-License-Identifier: Apache-2.0

package consistent

import (
	"fmt"
	"math/rand"
	"testing"
)

const (
	// objectSeed is the random number seed we use to create objects to hash
	objectSeed int64 = 7245298734452934458

	// objectCount is the number of random objects we generate for hash inputs
	objectCount int = 1000

	// serviceCount is the number of random service names to generate for hashes
	serviceCount int = 100
)

// hashObjects contains a standard set of random objects to hash to services.
var hashObjects [objectCount][16]byte

// services contains a number of service names to use in tests and benchmarks.
var services [serviceCount]string

func TestMain(m *testing.M) {
	random := rand.New(
		rand.NewSource(objectSeed),
	)

	for i := range len(hashObjects) {
		random.Read(hashObjects[i][:])
	}

	for i := range len(services) {
		services[i] = fmt.Sprintf("service-%d.example.net", i)
	}

	m.Run()
}
