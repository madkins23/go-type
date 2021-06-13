package datamap

import (
	"fmt"

	"github.com/mitchellh/mapstructure"
)

// Marshal converts an item into a nested map data structure.
// Currently depends on github.com/mitchellh/mapstructure.
func Marshal(item interface{}) (map[string]interface{}, error) {
	toMap := make(map[string]interface{})
	if err := mapstructure.Decode(item, &toMap); err != nil {
		return nil, fmt.Errorf("decode to map: %s", err)
	}

	return toMap, nil
}

// Unmarshal converts a nested map data structure into an item.
// Currently depends on github.com/mitchellh/mapstructure.
func Unmarshal(fromMap map[string]interface{}, item interface{}) error {
	if err := mapstructure.Decode(fromMap, item); err != nil {
		return fmt.Errorf("decode from map: %s", err)
	}

	return nil
}