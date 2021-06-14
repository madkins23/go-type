package serial

// There can be only one...

var theOne Mapper = NewMapper(nil)

func Highlander() Mapper {
	return theOne
}

func Marshal(item interface{}) (map[string]interface{}, error) {
	return theOne.Marshal(item)
}

func Unmarshal(data map[string]interface{}) (interface{}, error) {
	return theOne.Unmarshal(data)
}
