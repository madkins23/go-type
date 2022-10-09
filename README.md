# `go-type`

Go language type utilities.

See the [source](https://github.com/madkins23/go-type)
or [godoc](https://godoc.org/github.com/madkins23/go-type) for more detailed documentation.

[![Go Report Card](https://goreportcard.com/badge/github.com/madkins23/go-type)](https://goreportcard.com/report/github.com/madkins23/go-type)
![GitHub](https://img.shields.io/github/license/madkins23/go-type)
![GitHub release (latest by date)](https://img.shields.io/github/v/release/madkins23/go-type)
[![Go Reference](https://pkg.go.dev/badge/github.com/madkins23/go-type.svg)](https://pkg.go.dev/github.com/madkins23/go-type)

# Packages

Currently there is only one type utility that provides a registry for named Go types.
This is implemented in the `reg` package.
Other type utilities may be added in the future.

## Package `reg`

This package provides a dynamic type registry.
The Go language doesn't provide dynamic type lookup by type name.
All of the type names are removed during compilation,
so even using the `reflect` package they're invisible.

There is thus no way to create instances of types that are unknown to a package
but might be provided by an application using that package.
In cases requiring type instance creation by name (from a string)
it is necessary to provide a way to register types by name,
look them up, and then create new instances thereof.

The original motivation for this package is serialization of data
into JSON or YAML.
The deserialization of that data requires instance creation by type name.
See [`go-serial`](https://github.com/madkins23/go-serial)
for example usage of this package.