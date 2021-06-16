package json

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/madkins23/go-type/reg"
)

func MarshalFieldItem(field string, item interface{}, withComma bool, work io.Writer) error {
	if _, err := fmt.Fprint(work, "\"favorite\": "); err != nil {
		return fmt.Errorf("write open brace: %w", err)
	}

	if err := MarshalItem(item, withComma, work); err != nil {
		return fmt.Errorf("marshal item: %w", err)
	}

	return nil
}

func MarshalItem(item interface{}, withComma bool, work io.Writer) error {
	name, err := reg.NameFor(item)
	if err != nil {
		return fmt.Errorf("get type name for %#v: %w", item, err)
	}

	if _, err = fmt.Fprintf(work, "{\"%s\": \"%s\", \"value\": ", reg.TypeField, name); err != nil {
		return fmt.Errorf("write header: %w", err)
	}

	encoder := json.NewEncoder(work)
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(item); err != nil {
		return fmt.Errorf("encode item to JSON: %w", err)
	}

	if _, err = fmt.Fprintf(work, "}"); err != nil {
		return fmt.Errorf("write footer: %w", err)
	}

	if withComma {
		if _, err := fmt.Fprint(work, ","); err != nil {
			return fmt.Errorf("write comma: %w", err)
		}
	}

	return nil
}
