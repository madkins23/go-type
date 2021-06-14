package json

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/madkins23/go-type/reg"
	"github.com/madkins23/go-type/serial"

	"github.com/stretchr/testify/suite"

	"github.com/madkins23/go-type/test"
)

var _ serial.Mappable = &filmJson{}

// testMapper must be global to be accessible from filmJSON.
// Normally this would be provided by serial.Highlander(),
// but for testing purposes create it globally so it uses the local Registry object.
var testMapper serial.Mapper

type JsonTestSuite struct {
	suite.Suite
	film      *filmJson
	registry  reg.Registry
	mapper    serial.Mapper
	converter serial.Converter
}

func (suite *JsonTestSuite) SetupSuite() {
	suite.registry = reg.NewRegistry()
	suite.Assert().NotNil(suite.registry)
	suite.mapper = serial.NewMapper(suite.registry)
	suite.Assert().NotNil(suite.mapper)
	suite.Assert().NoError(suite.registry.AddAlias("test", test.Alpha{}), "creating test alias")
	suite.Assert().NoError(suite.registry.AddAlias("testJSON", &filmJson{}), "creating testJSON alias")
	suite.Assert().NoError(suite.registry.Register(&filmJson{}))
	suite.Assert().NoError(suite.registry.Register(&test.Alpha{}))
	suite.Assert().NoError(suite.registry.Register(&test.Bravo{}))
}

func (suite *JsonTestSuite) SetupTest() {
	suite.converter = NewConverter(suite.mapper)
	suite.Assert().NotNil(suite.converter)
	suite.film = &filmJson{Name: "Test JSON", Index: make(map[string]test.Actor)}
	suite.film.Lead = &test.Alpha{Name: "Goober", Percent: 13.23}
	suite.film.addActor("Goober", suite.film.Lead)
	suite.film.addActor("Snoofus", test.NewBravo(false, 17, "stuff"))
	suite.film.addActor("Noodle", test.NewAlpha("Noodle", 19.57, "stuff"))
	suite.film.addActor("Soup", &test.Bravo{Finished: true, Iterations: 79})
	testMapper = suite.mapper // no other convenient way to pass to where it's needed
}

func TestJsonSuite(t *testing.T) {
	suite.Run(t, new(JsonTestSuite))
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

//////////////////////////////////////////////////////////////////////////

func (suite *JsonTestSuite) TestConverterIsRegistry() {
	_, ok := suite.converter.(reg.Registry)
	suite.Assert().True(ok)
}

func (suite *JsonTestSuite) TestGetTypeName() {
	reader := strings.NewReader(simpleJson)
	suite.Assert().NotNil(reader)
	typeName, err := suite.converter.TypeName(reader)
	suite.Assert().NoError(err)
	suite.Assert().Equal("[testJSON]filmJson", typeName)
}

//////////////////////////////////////////////////////////////////////////

// TestMarshalCycle verifies the JSON Marshal/Unmarshal works as expected.
func (suite *JsonTestSuite) TestMarshalCycle() {
	bytes, err := json.Marshal(suite.film)
	suite.Assert().NoError(err)

	//fmt.Printf(">>> marshaled:\n%#v\n", string(bytes))

	var film filmJson
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

func (suite *JsonTestSuite) TestLoadFromString() {
	item, err := suite.converter.LoadFromString(simpleJson)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(item)
}

func (suite *JsonTestSuite) TestSaveToString() {
	text, err := suite.converter.SaveToString(suite.film)
	suite.Assert().NoError(err)
	suite.Assert().NotEmpty(text)
	suite.Assert().True(strings.Contains(text, `"<type>":"[testJSON]filmJson"`))
	suite.Assert().True(strings.Contains(text, `"<type>":"[test]Alpha"`))
	suite.Assert().True(strings.Contains(text, `"<type>":"[test]Bravo"`))
}

// TODO: Fix!
func (suite *JsonTestSuite) TestMarshalFileCycle() {
	file, err := ioutil.TempFile("", "*.test.json")
	suite.Assert().NoError(err)
	suite.Assert().NotNil(file)
	fileName := file.Name()
	// Go ahead and close it, just needed the file name.
	suite.Assert().NoError(file.Close())

	film := suite.film
	suite.Assert().NoError(suite.converter.SaveToFile(film, fileName))

	reloaded, err := suite.converter.LoadFromFile(fileName)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(reloaded)
	// TODO: Fix!
	//suite.Assert().Equal(film, reloaded)
}

func (suite *JsonTestSuite) TestMarshalStringCycle() {
	film := suite.film
	str, err := suite.converter.SaveToString(film)
	suite.Assert().NoError(err)
	suite.NotZero(str)

	fmt.Print("--- TestMarshalStringCycle ---------\n", str, "------------------------------------\n")

	reloaded, err := suite.converter.LoadFromString(str)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(reloaded)
	// TODO: Fix!
	//suite.Assert().Equal(film, reloaded)
}

//////////////////////////////////////////////////////////////////////////

type filmJson struct {
	json.Marshaler
	json.Unmarshaler

	Name  string
	Lead  test.Actor
	Cast  []test.Actor
	Index map[string]test.Actor
}

func (film *filmJson) addActor(name string, act test.Actor) {
	film.Cast = append(film.Cast, act)
	film.Index[name] = act
}

const fldName = "name"
const fldLead = "lead"
const fldCast = "cast"
const fldIndex = "index"

func (film *filmJson) Marshal() (map[string]interface{}, error) {
	var err error

	converted := make(map[string]interface{})

	converted[fldName] = film.Name

	if converted[fldLead], err = testMapper.Marshal(film.Lead); err != nil {
		return nil, fmt.Errorf("converting %s to map: %w", fldLead, err)
	}

	cast := make([]interface{}, len(film.Cast))
	for i, member := range film.Cast {
		if cast[i], err = testMapper.Marshal(member); err != nil {
			return nil, fmt.Errorf("converting %s member to map: %w", fldCast, err)
		}
	}
	converted[fldCast] = cast

	index := make(map[string]interface{}, len(film.Index))
	for key, member := range film.Index {
		if index[key], err = testMapper.Marshal(member); err != nil {
			return nil, fmt.Errorf("converting cast member to map: %w", err)
		}
	}
	converted[fldIndex] = index

	return converted, nil
}

func (film *filmJson) Unmarshal(mapData map[string]interface{}) error {
	var ok bool
	var err error
	if film.Name, ok = mapData[fldName].(string); !ok {
		return fmt.Errorf("bad film name")
	} else if film.Name == "" {
		return fmt.Errorf("no film name")
	}

	if film.Lead, err = film.unmarshalActor(mapData[fldLead]); err != nil {
		return fmt.Errorf("unmarshaling lead actor: %w", err)
	}

	if cast, ok := mapData[fldCast].([]interface{}); !ok {
		return fmt.Errorf("bad cast")
	} else {
		film.Cast = make([]test.Actor, len(cast))
		for i, member := range cast {
			if film.Cast[i], err = film.unmarshalActor(member); err != nil {
				return fmt.Errorf("unmarshaling cast member: %w", err)
			}
		}
	}

	if index := mapData[fldIndex].(map[string]interface{}); !ok {
		return fmt.Errorf("bad index")
	} else {
		film.Index = make(map[string]test.Actor, len(index))
		for name, member := range index {
			if film.Index[name], err = film.unmarshalActor(member); err != nil {
				return fmt.Errorf("unmarshaling index member: %w", err)
			}
		}
	}

	return nil
}

// MarshalJSON is called to marshal the filmJson object properly.
// It is necessary filmJson contains some interface fields that must be populated.
func (film *filmJson) MarshalJSON() ([]byte, error) {
	if toMap, err := film.Marshal(); err != nil {
		return nil, fmt.Errorf("pushing film to map: %w", err)
	} else {
		return json.Marshal(toMap)
	}
}

func (film *filmJson) UnmarshalJSON(input []byte) error {
	var err error
	mapData := make(map[string]interface{})
	if err = json.Unmarshal(input, &mapData); err != nil {
		return fmt.Errorf("unmarshal JSON to map data: %w", err)
	}

	return film.Unmarshal(mapData)
}

func (film *filmJson) unmarshalActor(input interface{}) (test.Actor, error) {
	actMap, ok := input.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("actor input should be map")
	} else if item, err := testMapper.Unmarshal(actMap); err != nil {
		return nil, fmt.Errorf("creating item from map: %w", err)
	} else if act, ok := item.(test.Actor); !ok {
		return nil, fmt.Errorf("item is not an actor")
	} else {
		return act, nil
	}
}

//////////////////////////////////////////////////////////////////////////

const simpleJson = `{
    "<type>": "[testJSON]filmJson",
    "name":   "Blockbuster Movie",
    "lead": {
        "<type>": "[test]Alpha",
        "name": "Lance Lucky",
        "percent": 23.79,
        "extra": "Yaaaa!"
    },
    "cast": [
        {
            "<type>": "[test]Alpha",
            "name": "Lance Lucky",
            "percent": 23.79,
            "extra": false
        },
        {
            "<type>": "[test]Bravo",
            "finished": true,
            "iterations": 13,
            "extra": "gibbering ghostwhistle"
        }
    ],
    "index": {
        "Lucky, Lance": {
            "<type>": "[test]Alpha",
            "name": "Lance Lucky",
            "percent": 23.79,
            "extra": "marshmallow stars"
        },
        "Queue, Susie": {
            "<type>": "[test]Bravo",
            "finished": true,
            "iterations": 13,
            "extra": 19.57
        }
    }
}
`
