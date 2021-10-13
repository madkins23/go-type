// Package reg provides a type registration mechanism to support instance creation by type name.
//
// The Reg package provides a dynamic type registry.
// The Go language doesn't provide dynamic type lookup by name.
// All type names are lost during compilation.
// In cases requiring creation of type instances by name (a string)
// it is necessary to provide a way to register types by name.
//
// Type Registration
//
// Creating a type registration mechanism provides a way to go from
// the type name to a usable instance of the type.
// Type registration uses reflection to convert an example instance
// of the type into a type name and type object during application initialization.
// After this the type object can be found by type name
// and reflection can be used to generate a new instance of the appropriate type.
//
// Types are registered early in program execution,
// generally in init() blocks or during static global variable initialization.
// During registration the name and type data is cached to speed reuse.
// Subsequently the registry can be used to look up object type by name,
// create new instance of type by name, and look up type name from an instance.
//
// Type Naming
//
// Full type names are acquired from the Go Type object.
// The full name is a combination of the package path
// (e.g. github.com/madkins23/go-type/reg) and type name.
//
// Registered types may have multiple names.
// The first one is the full type name, the rest are names created using aliases.
// When finding a name for a registered type via NameFor() the  shortest name will be returned.
// When finding the type from the name all names will be checked.
//
// Aliases
//
// Aliases may be defined for packages in order to reduce type name size
// when requesting a name via NameFor().
// An alias is specified with a string and an example instance from the package.
//
// After a reg.Registry is created use reg.Registry.Alias() to specify
// a name for the package, where the package is specified by a pointer
// to an example type from that package.
//
// The reg.Registry object also contains a map of aliases.
// When registering a type the full type name is acquired and stored.
// Since these can be long it is possible to provide an alias to a package
// which will be used in serialized objects.
//
// Global vs local Registry
//
// There is a global reg.Registry object created during initialization.
// The user may choose to use this via various functions
// or to create a local reg.Registry object and use its methods.
// The top-level functions that call the global reg.Registry object
// just use the methods.
//
// The problem with local reg.Registry objects is that they are not always available.
// Serializing objects provides a good example.
// The existing serialization libraries don't provide a way to
// attach data (in this case a reg.Registry object) to the encoder,
// nor do they pass context.Context objects down the call tree
// to be used by json.MarshalJSON() or yaml.MarshalYAML() or
// their unmarshal counterparts.
// In these cases a global reg.Registry is desirable if not necessary.
//
// While using global resources is generally considered bad,
// it is also good to consider why local registry objects might be needed.
// Is there some actual need to separate type registrations?
// After all, the types themselves are global.
//
// A single global reg.Registry object is provided via the reg.Highlander() function.
// In the general case this will be sufficient for all use.
//
// Concurrency
//
// The basic registry object is not guaranteed safe for concurrent access.
// Since the type registration should be done at application startup
// and subsequent access will be read-only to underlying map objects
// this is probably sufficient for most usage.
// If not, use reg.NewRegistrar() to create a Registry object that
// uses mutex locks.
package reg
