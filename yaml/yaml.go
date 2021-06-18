package yaml

import (
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/madkins23/go-type/reg"
)

type Wrapper struct {
	TypeName string
	Contents interface{}
}

func WrapItem(item interface{}) (*Wrapper, error) {
	w := &Wrapper{}
	return w, w.Wrap(item)
}

func (w *Wrapper) Wrap(item interface{}) error {
	var err error
	w.TypeName, err = reg.NameFor(item)
	if w.TypeName, err = reg.NameFor(item); err != nil {
		return fmt.Errorf("get type name for %#v: %w", item, err)
	}
	w.Contents = item

	return nil
}

func (w *Wrapper) Unwrap() (interface{}, error) {
	if w.TypeName == "" {
		return nil, fmt.Errorf("empty type field")
	} else if item, err := reg.Make(w.TypeName); err != nil {
		return nil, fmt.Errorf("make instance of type %s: %w", w.TypeName, err)
	} else {
		// TODO: fugly, how to fix?
		// Since the wrapper contents are map data, encode to YAML and then decode back.
		temp := &strings.Builder{}
		if err = yaml.NewEncoder(temp).Encode(w.Contents); err != nil {
			return nil, fmt.Errorf("encode contents to temp: %w", err)
		} else if err = yaml.NewDecoder(strings.NewReader(temp.String())).Decode(item); err != nil {
			return nil, fmt.Errorf("decode contents from temp: %w", err)
		} else {
			return item, nil
		}
	}
}
