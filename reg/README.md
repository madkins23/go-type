# `go-type/reg`

Type registration mechanism built with minimal use of reflection.

## Problem Description

Deserialization of generalized Go interfaces is problematic.
We want to generate *specific types* that implement the interface.
Most serialization mechanisms (e.g. JSON or YAML) don't  have any way of
storing the type information and existing serialization mechanisms only
create objects of the specific type provided to the method call.

### Manual Solution

We can write custom code to handle this problem.
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
Type registration uses reflection to convert the type name into
a type object during application initialization.
After this, reflection is used to generate a new instance of the
appropriate type when unmarshaling into the interface.

## Supported Formats

This package supports several serialization formats:

* JSON
* YAML

## Registration

Registration is done via a `reg.Registry` object.
This object stores information about a type by its full name.
Create a new object using `reg.NewRegistry()`.

The basic registry object is not safe for concurrent access.
Since the type registration should be done at application startup
this is probably sufficient.
If not, use `reg.NewRegistrar()` to create a Registry object that
uses mutex locks.

## Alias

The Registry object also contains a map of aliases.
When registering a type the full type name is provided.
The full type name is then stored in serialized objects.
Since these can be long it is possible to provide an alias to a package
which will be used in serialized objects.

After the Registry is created use `reg.Registry.Alias()` to specify
a name for the package, where the package is specified by a pointer
to an example type from that package.


