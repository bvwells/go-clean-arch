# go-clean-arch

[![GoDoc](http://godoc.org/github.com/bvwells/go-clean-arch?status.svg)](http://godoc.org/github.com/bvwells/go-clean-arch)
[![Build Status](https://travis-ci.org/bvwells/go-clean-arch.svg?branch=master)](https://travis-ci.org/bvwells/go-clean-arch)
[![Build status](https://ci.appveyor.com/api/projects/status/ea2u4hpy555b6ady?svg=true)](https://ci.appveyor.com/project/bvwells/go-clean-arch)
[![Go Report Card](https://goreportcard.com/badge/github.com/bvwells/go-clean-arch)](https://goreportcard.com/report/github.com/bvwells/go-clean-arch)

Linter for enforcing clean architecture principles in Go. The linter reports violations of clean architecture dependency rules by checking go imports 
within packages against user defined dependency rules.

For more information on 'The Clean Architecture Principle' see [Clean Architecture](https://8thlight.com/blog/uncle-bob/2012/08/13/the-clean-architecture.html) by Uncle Bob. The article [Applying The Clean Architecture To Go Applications](http://manuel.kiessling.net/2012/09/28/applying-the-clean-architecture-to-go-applications/) by Manual Kiessling is useful understand how clean architecture can be applied to Go applications.

## Installation

Install go-clean-arch with the following command:

```
go get -u github.com/bvwells/go-clean-arch
```

## Usage

To invoke go-clean arch use the command:

```
go-clean-arch [flags] [path ...]
```

where the flags are defined as:

    -c  Config file containing list of clean architecture layers from
        inner layers to outer laters.

## Example

The go-clean-arch linter can be run on the Git repo https://github.com/ManuelKiessling/go-cleanarchitecture by cloning the repo using the command;

```
git clone https://github.com/ManuelKiessling/go-cleanarchitecture
```

Run the linter with the command:

```
go-clean-arch -c layers.cfg path-to-repo\go-cleanarchitecture\src
```

where the layers config file contains the clean architecture layers:

```
domain
usecases
interfaces
infrastructure
```

## Go Versions Supported

The most recent major version of Go is supported. You can see which versions are
currently supported by looking at the lines following `go:` in
[`.travis.yml`](.travis.yml).