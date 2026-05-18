// SPDX-FileCopyrightText: 2026 Comcast Cable Communications Management, LLC
// SPDX-License-Identifier: Apache-2.0

package consistent

import (
	"fmt"
	"slices"

	"github.com/xmidt-org/medley"
)

func ExampleBuilder_simple() {
	hostNames := []string{
		"service-1.abc.something.net",
		"service-2.abc.something.net",
		"service-3.abc.something.net",
		"service-4.abc.something.net",
		"service-5.abc.something.net",
	}

	// use the default algorithm and vnodes
	builder := new(Builder[string]).
		ExpectedValues((len(hostNames)))

	// use medley.Objectify to use the host name itself as the hash object
	ring := builder.Build(
		medley.Objectify(
			medley.String,
			slices.Values(hostNames),
		),
	)

	// now we can assign clients to nodes
	hostName := ring.Nearest(medley.String("aclient"))
	fmt.Println(hostName)

	// Output:
	// service-1.abc.something.net
}

func ExampleBuilder_struct() {
	// you can hash to structs
	type service struct {
		hostName      string
		port          int
		favoriteThing string // doesn't matter ... you can have any fields you want
	}

	services := []*service{
		{hostName: "service-1.abc.something.net", port: 8080, favoriteThing: "pomegranates"},
		{hostName: "service-2.abc.something.net", port: 1111, favoriteThing: "apples"},
		{hostName: "service-3.abc.something.net", port: 54, favoriteThing: "watches"},
		{hostName: "service-4.abc.something.net", port: 8750, favoriteThing: "giraffes"},
		{hostName: "service-5.abc.something.net", port: 2562, favoriteThing: "aliens"},
	}

	// use a more interesting builder
	// you can use Builder[service] if you prefer
	builder := new(Builder[*service]).
		VNodes(10).
		Algorithm(medley.FNV64a()).
		ExpectedValues((len(services)))

	ring := builder.Build(
		medley.Objectify(
			func(s *service) medley.Object { return medley.String(s.hostName) },
			slices.Values(services),
		),
	)

	// now we can obtain a *service for a client
	svc := ring.Nearest(medley.String("aclient"))
	fmt.Println(svc.hostName)
	fmt.Println(svc.port)
	fmt.Println(svc.favoriteThing)

	// Output:
	// service-4.abc.something.net
	// 8750
	// giraffes
}
