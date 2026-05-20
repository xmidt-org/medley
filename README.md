# medley

medley provides service location based on hashing.

[![Build Status](https://github.com/xmidt-org/medley/actions/workflows/ci.yml/badge.svg)](https://github.com/xmidt-org/medley/actions/workflows/ci.yml)
[![codecov.io](http://codecov.io/github/xmidt-org/medley/coverage.svg?branch=main)](http://codecov.io/github/xmidt-org/medley?branch=main)
[![Go Report Card](https://goreportcard.com/badge/github.com/xmidt-org/medley)](https://goreportcard.com/report/github.com/xmidt-org/medley)
[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=xmidt-org_medley&metric=alert_status)](https://sonarcloud.io/dashboard?id=xmidt-org_medley)
[![Apache V2 License](http://img.shields.io/badge/license-Apache%20V2-blue.svg)](https://github.com/xmidt-org/medley/blob/main/LICENSE)
[![GitHub Release](https://img.shields.io/github/release/xmidt-org/medley.svg)](CHANGELOG.md)
[![GoDoc](https://pkg.go.dev/badge/github.com/xmidt-org/medley)](https://pkg.go.dev/github.com/xmidt-org/medley)

## Summary

Medley is a hashing package aimed at various types of distributed hashing. Currently, only consistent hashing is supported.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Install](#install)
- [Overview](#overview)
- [Key Features](#key-features)
  - [Generic medley.Hash interface](#generic-medleyhash-interface)
  - [Algorithm](#algorithm)
  - [Consistent Hashing](#consistent-hashing)
- [Contributing](#contributing)

## Code of Conduct

This project and everyone participating in it are governed by the [XMiDT Code Of Conduct](https://xmidt.io/code_of_conduct/). 
By participating, you agree to this Code.

## Install

go get -u github.com/xmidt-org/medley@latest

## Overview

Medley is a hashing library that provides some additional functionality on top of the `golang` stdlib package `hash`. Below are some key features of `medley`. See the [godoc](https://pkg.go.dev/github.com/xmidt-org/medley) for more information.

## Key Features

### Generic medley.Hash interface

Medley adds a [Hash](https://pkg.go.dev/github.com/xmidt-org/medley#Hash) interface that is genericized. This normalizes `hash.Hash32` and `hash.Hash64` in the standard library.

```golang
// 32-bit hashing
h32 := medley.AsHash32(fnv.New32a())
h32.WriteString("here is a string")
fmt.Println(h.Value()) // instead of Sum32()

// 64-bit hashing ... notice how similar it is
h64 := medley.AsHash32(fnv.New64a())
h64.WriteString("here is a string")
fmt.Println(h.Value()) // instead of Sum64()
```

The `WriteString` method computes the hash of a string without additional allocation.

### Algorithm

Medley exposes a generic [Algorithm](https://pkg.go.dev/github.com/xmidt-org/medley#Algorithm) type that is backed by a hashing package. Medley's `Default32` and `Default64` algorithms are based on [murmur3](https://pkg.go.dev/github.com/spaolacci/murmur3).

```golang
// An algorithm can create a hash
h64 := medley.Default64().New()
h64.Write([]byte{1, 2, 3})
fmt.Println(h64.Value())

// An Algorithm can also do one-time sums
medley.Default64().Sum([]byte{1, 2, 3})
medley.Default64().SumString("here is a string")
```

As with `Hash.WriteString`, `Algorithm.SumString` sums a string's bytes without additional allocation.

### Consistent Hashing

The [consistent](https://pkg.go.dev/github.com/xmidt-org/medley/consistent) package implements simple consistent hashing. A hash `Ring` is built using a `Builder`.

```golang
// First, create a ring using the desired configuration
builder := new(consistent.Builder[string, string]).
    Algorithm(medley.FNV64a()). // the default is medley.Default64()
    VNodes(10) // the default is consistent.DefaultVNodes

// Now, build rings based on sequences
hostNames := []string{"host-1.test.org", "host-2.test.org", "host-3.test.org"}
ring := builder.Build(
    len(hostNames), // optional. used as a hint for preallocation. can be zero or negative for no preallocation.
    medley.Stringify(slices.Values(hostNames)),
)

fmt.Println(ring.NearestString("myobject"))
```

Custom values are also supported on a `Ring`. You simply have to tell `medley` how to obtain the hashable object from each custom value.

```golang
type server struct {
    hostName string
    port int
}

servers := []*server{
    {hostName: "host-1.test.org", port: 1111},
    {hostName: "host-2.test.org", port: 2222},
    {hostName: "host-3.test.org", port: 2222},
}

var builder Builder[string, *server]

ring := builder.Build(
    len(servers),
    medley.Objectify(
        func(s *server) string { return s.hostName },
        slices.Values(servers),
    )
)

fmt.Println(ring.NearestString("myobject"))
```

## Contributing

Refer to [CONTRIBUTING.md](CONTRIBUTING.md).
