// Package reg provides a type registration mechanism.
// Types are registered by example with a Registry object
// which can then create new empty objects of the registered type
// and supports marshaling into/out of BSON, JSON, and YAML.
package reg

// TODO: modularize marshal/unmarshal support?
//  What does that even mean?
