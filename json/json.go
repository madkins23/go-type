package json

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/madkins23/go-type/reg"
	"github.com/madkins23/go-type/serial"
)

func NewConverter(mapper serial.Mapper) serial.Converter {
	return serial.NewConverter(&conversion{}, mapper)
}

//////////////////////////////////////////////////////////////////////////

type conversion struct {
}

var typeMatcher = regexp.MustCompile("\"" + reg.TypeFieldEscaped + "\":\\s*\"([^\"]+)\",")

func (c conversion) TypeName(reader io.ReadSeeker) (string, error) {
	buffered := bufio.NewReader(reader)

	for {
		if line, _, err := buffered.ReadLine(); err == io.EOF {
			break
		} else if err != nil {
			return "", fmt.Errorf("read line: %w", err)
		} else if matches := typeMatcher.FindStringSubmatch(string(line)); len(matches) < 1 {
			continue
		} else if _, err := reader.Seek(0, io.SeekStart); err != nil {
			return "", fmt.Errorf("reset reader: %w", err)
		} else {
			// Trim off any quotes and whitespace.
			return strings.Trim(matches[1], "'\" "), nil
		}
	}

	return "", fmt.Errorf("unable to locate type field")
}

func (c conversion) Decode(item interface{}, reader io.Reader) error {
	if err := json.NewDecoder(reader).Decode(item); err != nil {
		return fmt.Errorf("decode item from JSON: %w", err)
	}
	return nil
}

func (c conversion) Encode(item interface{}, writer io.Writer) error {
	encoder := json.NewEncoder(writer)
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(item); err != nil {
		return fmt.Errorf("encode item to JSON: %w", err)
	}
	return nil
}
