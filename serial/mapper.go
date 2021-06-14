package serial

import (
	"fmt"
	"reflect"

	"github.com/madkins23/go-type/reg"
)

type Mapper interface {
	reg.Registry
	Marshal(item interface{}) (map[string]interface{}, error)
	Unmarshal(data map[string]interface{}) (interface{}, error)
}

// NewMapper returns a new Mapper composed with the specified reg.Registry.
// If the specified reg.Registry is nil then reg.Highlander() is used to get the global registry.
func NewMapper(registry reg.Registry) Mapper {
	if registry == nil {
		registry = reg.Highlander()
	}

	return &mapper{Registry: registry}
}

// Mappable defines methods for marshaling an object to/from from a map.
type Mappable interface {
	Marshal() (map[string]interface{}, error)
	Unmarshal(map[string]interface{}) error
}

//////////////////////////////////////////////////////////////////////////

var _ Mapper = &mapper{}
var _ reg.Registry = &mapper{}

type mapper struct {
	reg.Registry
}

// Marshal converts an item of a registered type into a map for further processing.
// The registered type name will be put into the map in a special field named by TypeField.
// An error is returned if the type of the item is not registered.
// If the item implements Mappable then its Marshal method is called to populate the map.
func (mapper *mapper) Marshal(item interface{}) (map[string]interface{}, error) {
	value := reflect.ValueOf(item)
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}
	if value.Kind() != reflect.Struct {
		return nil, fmt.Errorf("item %s is not a struct", item)
	}

	var err error
	var result map[string]interface{}
	if regItem, ok := item.(Mappable); ok {
		if result, err = regItem.Marshal(); err != nil {
			return nil, fmt.Errorf("marshaling item to map: %w", err)
		}
	} else {
		result = make(map[string]interface{})
	}

	// Add the special marker for the type of the object.
	// This should work with both JSON and YAML.
	if result[reg.TypeField], err = mapper.NameFor(item); err != nil {
		return nil, fmt.Errorf("get name for %v: %w", value, err)
	}

	return result, nil
}

// Unmarshal attempts to return a new item of the type specified in the map.
// The registered type name is acquired from a special field named by TypeField in the map.
// An error is returned if this field does not exist or is not registered.
// If the item implements Mappable then its Unmarshal method is called to populate the map.
func (mapper *mapper) Unmarshal(data map[string]interface{}) (interface{}, error) {
	typeField, found := data[reg.TypeField]
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

	if mappable, ok := item.(Mappable); ok {
		if err := mappable.Unmarshal(data); err != nil {
			return nil, fmt.Errorf("pulling item fields from map: %w", err)
		}
	}

	return item, nil
}
