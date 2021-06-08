# go-type

Go language type registry.

Provides a registry for named Go classes.
The registry was constructed to support marshal/unmarshaling of Go objects.
Support is provided for BSON (binary JSON used by Mongo), JSON, and YAML.
This registry is intended to replace older code to be removed from
[`go-utils`](https://github.com/madkins23/go-utils).

You are more than welcome to use this package as is but these are
utility packages constructed by the author for use in personal projects.
The author makes occasional changes and attempts to follow proper versioning and release protocols,
however this code should not be considered production quality or maintained.

*Consider copying the code into your own project and modifying to fit your need.*

See the [source](https://github.com/madkins23/go-type)
or [godoc](https://godoc.org/github.com/madkins23/go-type) for documentation.

## `reg`

Type registration mechanism built with minimal use of reflection.

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
