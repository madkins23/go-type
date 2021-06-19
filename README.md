# `go-type`

Go language type utilities.

Currently there is only one type utility that provides a registry for named Go types.
This is implemented in the `reg` package.
Other type utilities may be added in the future.

### Caveats

You are more than welcome to use this software as is but these are
utility packages constructed by the author for use in personal projects.
The author makes occasional changes and attempts to follow proper versioning and release protocols,
however this code should not be considered production quality or maintained.

*Consider copying the code into your own project and modifying to fit your need.*

See the [source](https://github.com/madkins23/go-type)
or [godoc](https://godoc.org/github.com/madkins23/go-type) for documentation.

## Package `reg`

This package provides a dynamic type registry.
Since the Go language doesn't provide dynamic type lookup by name
there is no way to create instances of types that are unknown to a package
but might be provided by an application using that package.
In cases requiring type instance creation by name (from a string)
it is necessary to provide a way to register types by name,
look them up, and then create new instances thereof.

The original motivation for this package is serialization of data
into JSON or YAML.
The deserialization of that data requires instance creation by type name.
See [`go-serial`](https://github.com/madkins23/go-serial)
for example usage of this package.