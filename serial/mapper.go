package serial

import (
	"fmt"
	"reflect"

	"github.com/madkins23/go-type/reg"
)

type Mapper interface {
	reg.Registry
	ConvertItemToMap(item interface{}) (map[string]interface{}, error)
	CreateItemFromMap(in map[string]interface{}) (interface{}, error)
}

func NewMapper(registry reg.Registry) (Mapper, error) {
	if registry == nil {
		return nil, reg.ErrNilRegistry
	}

	return &mapper{Registry: registry}, nil
}

// Mappable defines methods for marshaling an object to/from from a map.
type Mappable interface {
	PushToMap(toMap map[string]interface{}) error
	PullFromMap(fromMap map[string]interface{}) error
}

//////////////////////////////////////////////////////////////////////////

type mapper struct {
	reg.Registry
}

// ConvertItemToMap converts an item of a registered type into a map for further processing.
// The registered type name will be put into the map in a special field named by TypeField.
// An error is returned if the type of the item is not registered.
// If the item implements Mappable then its PushToMap method is called to populate the map.
func (mapper *mapper) ConvertItemToMap(item interface{}) (map[string]interface{}, error) {
	value := reflect.ValueOf(item)
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}
	if value.Kind() != reflect.Struct {
		return nil, fmt.Errorf("item %s is not a struct", item)
	}

	result := make(map[string]interface{})

	// Add the special marker for the type of the object.
	// This should work with both JSON and YAML.
	var err error
	if result[reg.TypeField], err = mapper.NameFor(item); err != nil {
		return nil, fmt.Errorf("get name for %v: %w", value, err)
	}

	if regItem, ok := item.(Mappable); ok {
		if err := regItem.PushToMap(result); err != nil {
			return nil, fmt.Errorf("pushing item fields to map: %w", err)
		}
	}

	return result, nil
}

// CreateItemFromMap attempts to return a new item of the type specified in the map.
// The registered type name is acquired from a special field named by TypeField in the map.
// An error is returned if this field does not exist or is not registered.
// If the item implements Mappable then its PullFromMap method is called to populate the map.
func (mapper *mapper) CreateItemFromMap(in map[string]interface{}) (interface{}, error) {
	typeField, found := in[reg.TypeField]
	if !found {
		_ = fmt.Errorf("no object type in map")
	}
	typeName, ok := typeField.(string)
	if !ok {
		_ = fmt.Errorf("converting type field %v to string", typeField)
	}

	item, err := mapper.Make(typeName)
	if err != nil {
		return nil, fmt.Errorf("making item of type %s: %w", typeField, err)
	}

	if regItem, ok := item.(Mappable); ok {
		if err := regItem.PullFromMap(in); err != nil {
			return nil, fmt.Errorf("pulling item fields from map: %w", err)
		}
	}

	return item, nil
}
