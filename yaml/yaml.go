package yaml

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/madkins23/go-type/reg"
	"github.com/madkins23/go-type/serial"
)

func NewConverter(mapper serial.Mapper) (serial.Converter, error) {
	return serial.NewConverter(&conversion{}, mapper), nil
}

//////////////////////////////////////////////////////////////////////////

type conversion struct {
}

var typeMatcher = regexp.MustCompile("^" + reg.TypeFieldEscaped + ":\\s+(.+)$")

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
	if err := yaml.NewDecoder(reader).Decode(item); err != nil {
		return fmt.Errorf("unmarshal %w", err)
	}
	return nil
}

func (c conversion) Encode(item interface{}, writer io.Writer) error {
	encoder := yaml.NewEncoder(writer)
	if err := encoder.Encode(item); err != nil {
		_ = encoder.Close()
		return fmt.Errorf("marshal nexus: %w", err)
	}
	return encoder.Close()
}
