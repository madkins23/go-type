package convert

import (
	"fmt"

	"github.com/mitchellh/mapstructure"
)

// PushItemToMap provides a default mechanism for use in implementing Mappable.PushToMap().
// Currently depends on github.com/mitchellh/mapstructure.
func PushItemToMap(item interface{}, toMap map[string]interface{}) error {
	if err := mapstructure.Decode(item, &toMap); err != nil {
		return fmt.Errorf("decode to map: %s", err)
	}

	return nil
}

// PullItemFromMap provides a default mechanism for use in implementing Mappable.PullFromMap().
// Currently depends on github.com/mitchellh/mapstructure.
func PullItemFromMap(item interface{}, fromMap map[string]interface{}) error {
	if err := mapstructure.Decode(fromMap, item); err != nil {
		return fmt.Errorf("decode from map: %s", err)
	}

	return nil
}
