package serial

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/madkins23/go-type/reg"
)

type Converter interface {
	Conversion
	LoadFromFile(fileName string) (item interface{}, finalErr error)
	LoadFromString(source string) (interface{}, error)
	SaveToFile(item interface{}, fileName string) (finalErr error)
	SaveToString(item interface{}) (string, error)
}

// NewConverter returns a new Converter for a specific Conversion composed with the specified Mapper.
// If the specified Mapper is nil then Highlander() is used to get the global mapper.
func NewConverter(conversion Conversion, mapper Mapper) Converter {
	if conversion == nil {
		// I hate putting panic() anywhere just on general principles.
		// I also hate returning error results from "constructor" functions.
		// Since this is a really bad, nothing will work, kind of error I'm going with panic().
		panic("NewConverter() requires a conversion!")
	}

	if mapper == nil {
		mapper = Highlander()
	}

	return &converter{
		Conversion: conversion,
		Mapper:     mapper,
	}
}

//////////////////////////////////////////////////////////////////////////

var _ Converter = &converter{}
var _ Mapper = &converter{}
var _ reg.Registry = &converter{}

type converter struct {
	Conversion
	Mapper
}

// LoadFromFile loads an item of a registered type from the specified YAML file.
func (c *converter) LoadFromFile(fileName string) (item interface{}, finalErr error) {
	reader, err := os.Open(fileName)
	if err != nil {
		return nil, fmt.Errorf("open source file %s: %w", fileName, err)
	}
	defer func() {
		if err = reader.Close(); err != nil {
			finalErr = fmt.Errorf("close file: %w", err)
		}
	}()

	if item, err = c.loadFrom(reader); finalErr != nil {
		return nil, fmt.Errorf("load from readSeeker: %w", err)
	} else {
		return item, nil
	}
}

// LoadFromString loads an item of a registered type from a YAML string.
func (c *converter) LoadFromString(source string) (interface{}, error) {
	if item, err := c.loadFrom(strings.NewReader(source)); err != nil {
		return nil, fmt.Errorf("load from readSeeker: %w", err)
	} else {
		return item, nil
	}
}

// loadFrom is the common item creation/decode method used by other LoadXXX methods.
// The reader must be io.ReadSeeker to enable
// the stream to be reset after acquiring the type name if necessary.
// The embedded Conversion object is used for decoding to the created item.
func (c *converter) loadFrom(reader io.ReadSeeker) (interface{}, error) {
	if typeName, err := c.TypeName(reader); err != nil {
		return nil, fmt.Errorf("get type name: %w", err)
	} else if item, err := c.Make(typeName); err != nil {
		return nil, fmt.Errorf("make item of type %s: %w", typeName, err)
	} else if err := c.Decode(item, reader); err != nil {
		return nil, fmt.Errorf("decode %s item: %w", typeName, err)
	} else {
		return item, nil
	}
}

//////////////////////////////////////////////////////////////////////////

// SaveToFile saves an item of a registered type to the specified YAML file.
// TODO: The file is always created, perhaps other options should be considered?
func (c *converter) SaveToFile(item interface{}, fileName string) (finalErr error) {
	file, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("creating output file '%s': %w", fileName, err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			finalErr = fmt.Errorf("error closing source file: %w", err)
		}
	}()

	if err := c.saveTo(item, file); err != nil {
		return fmt.Errorf("save to writer: %w", err)
	} else {
		return nil
	}
}

// SaveToString marshals an item of a registered type to a YAML string.
func (c *converter) SaveToString(item interface{}) (string, error) {
	builder := &strings.Builder{}

	if err := c.saveTo(item, builder); err != nil {
		return "", fmt.Errorf("save to writer: %w", err)
	} else {
		return builder.String(), nil
	}
}

// saveTo is the common marshal/encode method used by other SaveXXX methods.
// The embedded Conversion object is used for encoding the marshaled data map.
func (c *converter) saveTo(item interface{}, writer io.Writer) error {
	// Marshaling the item via the converter adds the type field to the item.
	if dataMap, err := c.Marshal(item); err != nil {
		return fmt.Errorf("marshal item to map %w", err)
	} else if err := c.Encode(dataMap, writer); err != nil {
		return fmt.Errorf("encode data map: %w", err)
	}
	return nil
}
