# medley

medley does something good.

[![Build Status](https://github.com/xmidt-org/medley/actions/workflows/ci.yml/badge.svg)](https://github.com/xmidt-org/medley/actions/workflows/ci.yml)
[![codecov.io](http://codecov.io/github/xmidt-org/medley/coverage.svg?branch=main)](http://codecov.io/github/xmidt-org/medley?branch=main)
[![Go Report Card](https://goreportcard.com/badge/github.com/xmidt-org/medley)](https://goreportcard.com/report/github.com/xmidt-org/medley)
[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=xmidt-org_medley&metric=alert_status)](https://sonarcloud.io/dashboard?id=xmidt-org_medley)
[![Apache V2 License](http://img.shields.io/badge/license-Apache%20V2-blue.svg)](https://github.com/xmidt-org/medley/blob/main/LICENSE)
[![GitHub Release](https://img.shields.io/github/release/xmidt-org/medley.svg)](CHANGELOG.md)
[![GoDoc](https://pkg.go.dev/badge/github.com/xmidt-org/medley)](https://pkg.go.dev/github.com/xmidt-org/medley)

## Setup

1. Search and replace medley with your project name.
1. Initialize `go.mod` file: `go mod init github.com/xmidt-org/medley`
1. Add org teams to project (Settings > Manage Access): 
    - xmidt-org/admins with Admin role
    - xmidt-org/server-writers with Write role
1. Manually create the first release.  After v0.0.1 exists, other releases will be made by automation after the CHANGELOG is updated to reflect a new version header and nothing under the Unreleased header.
1. For libraries:
    1. Add org workflows in dir `.github/workflows`: push, tag, and release. This can be done by going to the Actions tab for the repo on the github site.
    1. Remove the following files/dirs: `.dockerignore`, `Dockerfile`, `Makefile`, `rpkg.macros`, `medley.yaml`, `deploy/`, and `conf/`.


## Summary

Medley is a consistent hash package that also exposes a simple API for creating additional hash strategies.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Details](#details)
- [Install](#install)
- [Contributing](#contributing)

## Code of Conduct

This project and everyone participating in it are governed by the [XMiDT Code Of Conduct](https://xmidt.io/code_of_conduct/). 
By participating, you agree to this Code.

## Contributing

Refer to [CONTRIBUTING.md](CONTRIBUTING.md).
