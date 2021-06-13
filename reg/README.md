# `go-type/reg`

Type registration mechanism to support instance creation by class name.

## Problem Description

Deserialization of generalized Go interfaces is problematic.
The applicadtion must generate *specific types* to implement an interface.
Most serialization mechanisms (e.g. JSON or YAML) don't  have any way of
storing the type information and existing serialization mechanisms only
create objects of the specific type provided to the method call
and so can't handle deserializing to interface fields.

### Manual Solution

It is possible to write custom code to handle this problem.
Generally it involves:

1. unmarshal to `interface{}` producing `map[string]interface{}`
2. use information in the `map`
   (e.g. a `Type` string field generated during serialization)
   to determine the type of the object
3. unmarshal from the map into a pointer of the correct type

There are various examples of this available online.

This is a viable strategy, but step 3 requires code from each type
to be embedded into the `UnmarshalXXX` method for any struct or array
containing the interface in question (interfaces can't have methods),
resulting in a potential code maintenance issue.

### Type Registration

Creating a type registration mechanism provides a way to go from
the type name to a usable instance of the type
without hard-coding the types into the unmarshal code.
Type registration uses reflection to convert an example instance
of the class into a type name and type object during application initialization.
After this the type object can be found by type name
and reflection can be used to generate a new instance of the appropriate type.

Classes are registered early in program execution,
generally in `init` blocks or during static global variable initialization.
During registration the name and type data is cached to speed reuse.
Subsequently the registry can be used to:

* look up type object by name
* create new instance of type by name
* look up type name from the instance

### Type Naming

Full type names are acquired from the Go `Type` object.
The full name is a combination of the package path
(e.g. `github.com/madkins23/go-type/reg`) and type name.

Registered types may have multiple names.
The first one is the full type name, the rest are names
created using aliases.
When finding a name for a registered type (e.g. during serialization)
the  shortest name will be returned.
When finding the type from the name all names will be checked.

### Aliases

Aliases may be defined for packages in order to reduce type name size
in serialized data.
An alias is specified with a string and an example instance from the package.

After the Registry is created use `reg.Registry.Alias()` to specify
a name for the package, where the package is specified by a pointer
to an example type from that package.

The Registry object also contains a map of aliases.
When registering a type the full type name is provided.
The full type name is then stored in serialized objects.
Since these can be long it is possible to provide an alias to a package
which will be used in serialized objects.

## Supported Formats

This package supports several serialization formats:

* BSON (binary JSON, used in Mongo DB)
* JSON
* YAML

## Registration

Registration is done via a `reg.Registry` object.
This object stores information about a type by its full name.
Create a new object using `reg.NewRegistry()`.

The basic registry object is not guaranteed safe for concurrent access.
Since the type registration should be done at application startup
and subsequent access will be read-only to underlying `map` objects
this is probably sufficient for most usage.
If not, use `reg.NewRegistrar()` to create a Registry object that
uses mutex locks.

A single global `reg.Registry` object is provided via the `reg.Highlander()` function.
In the general case this will be sufficient for all use.
