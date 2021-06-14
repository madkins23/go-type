package serial

// There can be only one...

// theOne is the global Mapper object returned by Highlander().
var theOne = NewMapper(nil)

// Highlander returns a global Mapper object created during initialization.
// This Mapper object uses reg.Highlander() for its Registry object.
// Normally there will only be one Mapper in use for the entire application.
// It is not necessary to use the global Mapper, it is just convenient.
func Highlander() Mapper {
	return theOne
}

// Marshal invokes serial.Highlander().Marshal().
func Marshal(item interface{}) (map[string]interface{}, error) {
	return theOne.Marshal(item)
}

// Unmarshal invokes serial.Highlander().Unmarshal().
func Unmarshal(data map[string]interface{}) (interface{}, error) {
	return theOne.Unmarshal(data)
}
