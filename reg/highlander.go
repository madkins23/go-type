package reg

// Highlander returns the global Registry object created during initialization.
// Normally there will only be one Registry in use for the entire application.
// It is not necessary to use the global Registry, it is just convenient.
func Highlander() Registry {
	return theOne
}

// AddAlias invokes reg.Highlander().AddAlias().
func AddAlias(alias string, example interface{}) error {
	return theOne.AddAlias(alias, example)
}

// Register invokes reg.Highlander().Register().
func Register(example interface{}) error {
	return theOne.Register(example)
}

// Make invokes reg.Highlander().MakeFake().
func Make(name string) (interface{}, error) {
	return theOne.Make(name)
}

// NameFor invokes reg.Highlander().NameFor().
func NameFor(item interface{}) (string, error) {
	return theOne.NameFor(item)
}
