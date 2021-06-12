package reg

// There can be only one...

var theOne Registry = NewRegistry()

func AddAlias(alias string, example interface{}) error {
	return theOne.AddAlias(alias, example)
}

func Register(example interface{}) error {
	return theOne.Register(example)
}

func Make(name string) (interface{}, error) {
	return theOne.Make(name)
}

func NameFor(item interface{}) (string, error) {
	return theOne.NameFor(item)
}
