# `go-type`

Go language type utilities.

Currently there is only one type utility which provides a registry for named Go classes.
This is implemented in the `reg` package.
Other type utilities may be added in the future.

You are more than welcome to use this software as is but these are
utility packages constructed by the author for use in personal projects.
The author makes occasional changes and attempts to follow proper versioning and release protocols,
however this code should not be considered production quality or maintained.

*Consider copying the code into your own project and modifying to fit your need.*

See the [source](https://github.com/madkins23/go-type)
or [godoc](https://godoc.org/github.com/madkins23/go-type) for documentation.

## `reg`

This package provides a dynamic type registry.
Since the Go language doesn't provide dynamic class lookup by name
it is necessary to provide a way to register classes by name
so that instances of those classes can be created by name later.

This registry was mainly constructed to support marshal/unmarshaling of Go objects.
Support is provided for BSON (binary JSON used by Mongo), JSON, and YAML.

* `reg.Alias` creates an aliased Registry object for package-specific registration.
  First registration via this object does the alias registration automatically.

* `reg.Registry` provides a way to register types by name.
  Normally Go doesn't keep type names at runtime, so it must be done by the application.
  The `Registry` object provides a way to track this and to generate objects of a "named" type.
  Created for use in Marshaling/Unmarshaling objects. Uses reflection.
  Not thread-safe but in normal usage pattern may not matter.
  See test files for examples of JSON and YAML marshal/unmarshal.

* `reg.Registrar` provides a thread-safe `Registry`.
  `Registry` methods are wrapped with a mutex object.
