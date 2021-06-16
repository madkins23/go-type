package datamap

import (
	"fmt"

	"github.com/mitchellh/mapstructure"
)

type DataMap map[string]interface{}

// Marshal converts an item into a nested map data structure.
// Currently depends on github.com/mitchellh/mapstructure.
func Marshal(item interface{}) (DataMap, error) {
	toMap := make(DataMap)
	if err := mapstructure.Decode(item, &toMap); err != nil {
		return nil, fmt.Errorf("decode to map: %s", err)
	}

	return toMap, nil
}

// Unmarshal converts a nested map data structure into an existing item.
// The item must be provided instead of returned as Unmarshal() won't otherwise know what type to use.
// In addition this function can be used to populate existing items from within their own methods.
// Currently depends on github.com/mitchellh/mapstructure.
func Unmarshal(fromMap DataMap, item interface{}) error {
	if err := mapstructure.Decode(fromMap, &item); err != nil {
		return fmt.Errorf("decode from map: %s", err)
	}

	return nil
}
