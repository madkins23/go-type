package reg

var singleton = NewRegistry()

// Singleton returns the global Registry object created during initialization.
// Normally there will only be one Registry in use for the entire application.
// It is not necessary to use the global Registry, it is just convenient.
func Singleton() Registry {
	return singleton
}

// SetSingleton sets the global Registry object to the specified registry.
func SetSingleton(registry Registry) {
	singleton = registry
}

// =============================================================================
// Apologies for the (now deprecated) Highlander nomenclature.  ;-)
// There can be only one...

// Highlander returns the global Registry object created during initialization.
// Normally there will only be one Registry in use for the entire application.
// It is not necessary to use the global Registry, it is just convenient.
//
// Deprecated: use Singleton instead, it is a more correct name.
func Highlander() Registry {
	return singleton
}

// Quicken sets the global Registry object to the specified registry.
//
// Deprecated: use SetSingleton instead, it is a more correct name.
func Quicken(registry Registry) {
	singleton = registry
}

// =============================================================================

// AddAlias invokes reg.Singleton().AddAlias().
func AddAlias(alias string, example interface{}) error {
	return singleton.AddAlias(alias, example)
}

// Make invokes reg.Singleton().Make().
func Make(name string) (interface{}, error) {
	return singleton.Make(name)
}

// NameFor invokes reg.Singleton().NameFor().
func NameFor(item interface{}) (string, error) {
	return singleton.NameFor(item)
}

// Register invokes reg.Singleton().Register().
func Register(example interface{}) error {
	return singleton.Register(example)
}
