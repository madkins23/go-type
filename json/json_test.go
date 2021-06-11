package json

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/madkins23/go-type/reg"
	"github.com/madkins23/go-type/serial"

	"github.com/stretchr/testify/suite"

	"github.com/madkins23/go-type/test"
)

// These tests demonstrate and validate use of a Registry to marshal/unmarshal JSON.

type JsonTestSuite struct {
	suite.Suite
	film     *filmJson
	registry reg.Registry
	mapper   serial.Mapper
}

func (suite *JsonTestSuite) SetupSuite() {
	var err error
	suite.registry = reg.NewRegistry()
	suite.Assert().NotNil(suite.registry)
	suite.mapper, err = serial.NewMapper(suite.registry)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(suite.mapper)
	suite.Assert().NoError(suite.registry.Alias("test", test.Alpha{}), "creating test alias")
	suite.Assert().NoError(suite.registry.Register(&test.Alpha{}))
	suite.Assert().NoError(suite.registry.Register(&test.Bravo{}))
}

func (suite *JsonTestSuite) SetupTest() {
	suite.film = &filmJson{Name: "Test JSON", Index: make(map[string]test.Actor)}
	suite.film.Lead = &test.Alpha{Name: "Goober", Percent: 13.23}
	suite.film.addActor("Goober", suite.film.Lead)
	suite.film.addActor("Snoofus", test.NewBravo(false, 17, "stuff"))
	suite.film.addActor("Noodle", test.NewAlpha("Noodle", 19.57, "stuff"))
	suite.film.addActor("Soup", &test.Bravo{Finished: true, Iterations: 79})
}

func TestJsonSuite(t *testing.T) {
	suite.Run(t, new(JsonTestSuite))
}

//////////////////////////////////////////////////////////////////////////

type filmJson struct {
	serial.WithMapper
	json.Marshaler
	json.Unmarshaler

	Name  string
	Lead  test.Actor
	Cast  []test.Actor
	Index map[string]test.Actor
}

type filmJsonConvert struct {
	Name  string
	Lead  interface{}
	Cast  []interface{}
	Index map[string]interface{}
}

func (film *filmJson) addActor(name string, act test.Actor) {
	film.Cast = append(film.Cast, act)
	film.Index[name] = act
}

func (film *filmJson) MarshalJSON() ([]byte, error) {
	var err error

	convert := filmJsonConvert{
		Name: film.Name,
	}

	if convert.Lead, err = film.Mapper().ConvertItemToMap(film.Lead); err != nil {
		return nil, fmt.Errorf("converting lead to map: %w", err)
	}

	convert.Cast = make([]interface{}, len(film.Cast))
	for i, member := range film.Cast {
		if convert.Cast[i], err = film.Mapper().ConvertItemToMap(member); err != nil {
			return nil, fmt.Errorf("converting cast member to map: %w", err)
		}
	}

	convert.Index = make(map[string]interface{}, len(film.Index))
	for key, member := range film.Index {
		if convert.Index[key], err = film.Mapper().ConvertItemToMap(member); err != nil {
			return nil, fmt.Errorf("converting cast member to map: %w", err)
		}
	}

	return json.Marshal(convert)
}

func (film *filmJson) UnmarshalJSON(input []byte) error {
	var err error

	convert := filmJsonConvert{}
	if err = json.Unmarshal(input, &convert); err != nil {
		return fmt.Errorf("unmarshaling input JSON into struct: %w", err)
	}

	film.Name = convert.Name

	if film.Lead, err = film.unmarshalActor(convert.Lead); err != nil {
		return fmt.Errorf("unmarshaling lead actor: %w", err)
	}

	film.Cast = make([]test.Actor, len(convert.Cast))
	for i, member := range convert.Cast {
		if film.Cast[i], err = film.unmarshalActor(member); err != nil {
			return fmt.Errorf("unmarshaling cast member: %w", err)
		}
	}

	film.Index = make(map[string]test.Actor, len(convert.Index))
	for name, member := range convert.Index {
		if film.Index[name], err = film.unmarshalActor(member); err != nil {
			return fmt.Errorf("unmarshaling index member: %w", err)
		}
	}

	return nil
}

func (film *filmJson) unmarshalActor(input interface{}) (test.Actor, error) {
	actMap, ok := input.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("actor input should be map")
	} else if mapper := film.Mapper(); mapper == nil {
		return nil, fmt.Errorf("no mapper on filmJson")
	} else if item, err := mapper.CreateItemFromMap(actMap); err != nil {
		return nil, fmt.Errorf("creating item from map: %w", err)
	} else if act, ok := item.(test.Actor); !ok {
		return nil, fmt.Errorf("item is not an actor")
	} else {
		return act, nil
	}
}

func copyItemToMap(toMap map[string]interface{}, fromItem interface{}) error {
	if bytes, err := json.Marshal(fromItem); err != nil {
		return fmt.Errorf("marshaling from %v: %w", fromItem, err)
	} else if err = json.Unmarshal(bytes, &toMap); err != nil {
		return fmt.Errorf("unmarshaling to %v: %w", toMap, err)
	}

	return nil
}

func copyMapToItem(toItem interface{}, fromMap map[string]interface{}) error {
	if bytes, err := json.Marshal(fromMap); err != nil {
		return fmt.Errorf("marshaling from %v: %w", fromMap, err)
	} else if err = json.Unmarshal(bytes, toItem); err != nil {
		return fmt.Errorf("unmarshaling to %v: %w", toItem, err)
	}

	return nil
}

//////////////////////////////////////////////////////////////////////////

// TestExample duplicates the YAML test.
// TODO: Not directly applicable to this test suite.
func (suite *JsonTestSuite) TestExample() {
	type T struct {
		F int `json:"a,omitempty"`
		B int
	}
	t := T{F: 1, B: 2}
	bytes, err := json.Marshal(t)
	suite.Assert().NoError(err)
	var x T
	suite.Assert().NoError(json.Unmarshal(bytes, &x))
	suite.Assert().Equal(t, x)
}

func (suite *JsonTestSuite) TestCycle() {
	suite.film.Open(suite.mapper)
	defer suite.film.Close()

	bytes, err := json.Marshal(suite.film)
	suite.Assert().NoError(err)

	//fmt.Printf(">>> marshaled:\n%s\n", string(bytes))

	var film filmJson
	film.Open(suite.mapper)
	defer film.Close()

	suite.Assert().NoError(json.Unmarshal(bytes, &film))
	suite.Assert().NotEqual(suite.film, &film) // fails due to unexported field 'extra'
	for _, act := range suite.film.Cast {
		// Remove unexported field.
		if a, ok := act.(*test.Alpha); ok {
			a.ClearExtra()
		} else if b, ok := act.(*test.Bravo); ok {
			b.ClearExtra()
		}
	}
	suite.Assert().Equal(suite.film, &film) // succeeds now that unexported fields are gone.
}