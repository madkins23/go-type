package reg

import (
	"fmt"

	"github.com/mitchellh/mapstructure"
)

// Mappable defines methods for marshaling an object to/from from a map.
type Mappable interface {
	PushToMap(toMap map[string]interface{}) error
	PullFromMap(fromMap map[string]interface{}) error
}

// PushMappableToMap provides a default mechanism for use in PushToMap().
// Currently depends on github.com/mitchellh/mapstructure.
func PushMappableToMap(mapper Mappable, toMap map[string]interface{}) error {
	if err := mapstructure.Decode(mapper, &toMap); err != nil {
		return fmt.Errorf("decode to map: %s", err)
	}

	return nil
}

// PullMappableFromMap provides a default mechanism for use in PullFromMap().
// Currently depends on github.com/mitchellh/mapstructure.
func PullMappableFromMap(mapper Mappable, fromMap map[string]interface{}) error {
	if err := mapstructure.Decode(fromMap, mapper); err != nil {
		return fmt.Errorf("decode from map: %s", err)
	}

	return nil
}
