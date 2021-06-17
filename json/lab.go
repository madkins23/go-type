package json

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/madkins23/go-type/data"

	"github.com/madkins23/go-type/reg"
)

type Wrapper struct {
	TypeName string
	Contents json.RawMessage
}

func WrapItem(item interface{}) (*Wrapper, error) {
	name, err := reg.NameFor(item)
	if err != nil {
		return nil, fmt.Errorf("get type name for %#v: %w", item, err)
	}

	build := &strings.Builder{}
	encoder := json.NewEncoder(build)
	encoder.SetEscapeHTML(false)
	if err = encoder.Encode(item); err != nil {
		return nil, fmt.Errorf("marshal wrapper contents: %w", err)
	}

	return &Wrapper{
		TypeName: name,
		Contents: json.RawMessage(build.String()),
	}, nil
}

func UnwrapItem(dataMap data.Map) (interface{}, error) {
	if fieldValue, found := dataMap["TypeName"]; !found {
		return nil, fmt.Errorf("no type field")
	} else if typeName, ok := fieldValue.(string); !ok {
		return nil, fmt.Errorf("bad type name: %#v", typeName)
	} else if typeName == "" {
		return nil, fmt.Errorf("empty type field")
	} else if fieldValue, found = dataMap["Contents"]; !found {
		return nil, fmt.Errorf("no content field")
	} else if item, err := reg.Make(typeName); err != nil {
		return nil, fmt.Errorf("make instance of type %s: %w", typeName, err)
	} else {
		// TODO: Irritating:
		// Re-construct contents into JSON so it can be decoded again.
		build := &strings.Builder{}
		encoder := json.NewEncoder(build)
		if err = encoder.Encode(fieldValue); err != nil {
			return nil, fmt.Errorf("marshal data for decoding: %w", err)
		} else if err = json.NewDecoder(strings.NewReader(build.String())).Decode(item); err != nil {
			return nil, fmt.Errorf("decode wrapper contents: %w", err)
		} else {
			return item, nil
		}
	}
}
