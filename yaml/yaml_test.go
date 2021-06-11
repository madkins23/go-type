package yaml

import (
	"fmt"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/madkins23/go-type/serial"

	"github.com/stretchr/testify/suite"
	"gopkg.in/yaml.v3"

	"github.com/madkins23/go-type/reg"
	"github.com/madkins23/go-type/test"
)

// These tests demonstrates and validates use of a Registry to marshal/unmarshal YAML.

type YamlTestSuite struct {
	suite.Suite
	film     *filmYaml
	registry reg.Registry
	mapper   serial.Mapper
}

func (suite *YamlTestSuite) SetupSuite() {
	var err error
	suite.registry = reg.NewRegistry()
	suite.Assert().NotNil(suite.registry)
	suite.mapper, err = serial.NewMapper(suite.registry)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(suite.mapper)
	suite.Assert().NoError(suite.registry.Alias("test", test.Alpha{}), "creating test alias")
	suite.Assert().NoError(suite.registry.Alias("testYAML", &filmYaml{}), "creating testYAML alias")
	suite.Assert().NoError(suite.registry.Register(&filmYaml{}))
	suite.Assert().NoError(suite.registry.Register(&test.Alpha{}))
	suite.Assert().NoError(suite.registry.Register(&test.Bravo{}))
}

func (suite *YamlTestSuite) SetupTest() {
	//test.CopyMapFromItemFn = func(toMap map[string]interface{}, fromItem interface{}) error {
	//	return copyViaYaml(toMap, fromItem)
	//}
	//test.CopyItemFromMapFn = func(toItem interface{}, fromMap map[string]interface{}) error {
	//	return copyViaYaml(toItem, fromMap)
	//}

	suite.film = &filmYaml{Name: "Test YAML", Index: make(map[string]test.Actor)}
	suite.film.Lead = &test.Alpha{Name: "Goober", Percent: 13.23}
	suite.film.addActor("Goober", suite.film.Lead)
	suite.film.addActor("Snoofus", test.NewBravo(false, 17, "stuff"))
	suite.film.addActor("Noodle", test.NewAlpha("Noodle", 19.57, "stuff"))
	suite.film.addActor("Soup", test.NewBravo(true, 79, ""))
}

func TestYamlSuite(t *testing.T) {
	suite.Run(t, new(YamlTestSuite))
}

//////////////////////////////////////////////////////////////////////////

func (suite *YamlTestSuite) TestConverterIsRegistry() {
	var conv interface{}
	var err error
	conv, err = NewConverter(suite.mapper)
	suite.Assert().NoError(err)
	_, ok := conv.(serial.Mapper)
	suite.Assert().True(ok)
}

func (suite *YamlTestSuite) TestGetTypeName() {
	reader := strings.NewReader(simpleYaml)
	suite.Assert().NotNil(reader)

	// Get type name.
	name, err := getYamlTypeName(reader)
	suite.Assert().NoError(err)
	suite.Assert().Equal("[testYAML]filmYaml", name)
}

//////////////////////////////////////////////////////////////////////////

type filmYaml struct {
	serial.WithMapper
	yaml.Marshaler
	yaml.Unmarshaler

	Name  string
	Lead  test.Actor
	Cast  []test.Actor
	Index map[string]test.Actor
}

func (film *filmYaml) addActor(name string, act test.Actor) {
	film.Cast = append(film.Cast, act)
	film.Index[name] = act
}

func (film *filmYaml) PullFromMap(fromMap map[string]interface{}) error {
	var ok bool
	if fromMap["name"] != nil {
		if film.Name, ok = fromMap["name"].(string); !ok {
			return fmt.Errorf("film name is not a string")
		}
	}

	var err error
	if fromMap["lead"] != nil {
		if film.Lead, err = film.pullActorFromMap(fromMap["lead"]); err != nil {
			return fmt.Errorf("pull lead actor from map: %w", err)
		}
	}

	if castElement, found := fromMap["cast"]; found && castElement != nil {
		if cast, ok := castElement.([]interface{}); ok {
			film.Cast = make([]test.Actor, 0, len(cast))
			for _, actMap := range cast {
				if act, err := film.pullActorFromMap(actMap); err != nil {
					return fmt.Errorf("pulling actor from map: %w", err)
				} else {
					film.Cast = append(film.Cast, act)
				}
			}
		}
	}

	if indexElement, found := fromMap["index"]; found && indexElement != nil {
		if index, ok := indexElement.(map[string]interface{}); ok {
			film.Index = make(map[string]test.Actor)
			for key, actMap := range index {
				if act, err := film.pullActorFromMap(actMap); err != nil {
					return fmt.Errorf("pulling actor from map: %w", err)
				} else {
					film.Index[key] = act
				}
			}
		}
	}

	return nil
}

func (film *filmYaml) pullActorFromMap(from interface{}) (test.Actor, error) {
	if fromMap, ok := from.(map[string]interface{}); !ok {
		return nil, fmt.Errorf("from is not a map")
	} else if mapper := film.Mapper(); mapper == nil {
		return nil, fmt.Errorf("no mapper on filmYaml")
	} else if actItem, err := mapper.CreateItemFromMap(fromMap); err != nil {
		return nil, fmt.Errorf("create actor from map: %w", err)
	} else if act, ok := actItem.(test.Actor); !ok {
		return nil, fmt.Errorf("created item is not an actor")
	} else {
		return act, nil
	}
}

func (film *filmYaml) PushToMap(toMap map[string]interface{}) error {
	var err error

	toMap["name"] = film.Name

	if toMap["lead"], err = film.Mapper().ConvertItemToMap(film.Lead); err != nil {
		return fmt.Errorf("converting lead to map: %w", err)
	}

	cast := make([]interface{}, len(film.Cast))
	for i, member := range film.Cast {
		if cast[i], err = film.Mapper().ConvertItemToMap(member); err != nil {
			return fmt.Errorf("converting cast member to map: %w", err)
		}
	}
	toMap["cast"] = cast

	index := make(map[string]interface{}, len(film.Index))
	for key, member := range film.Index {
		if index[key], err = film.Mapper().ConvertItemToMap(member); err != nil {
			return fmt.Errorf("converting cast member to map: %w", err)
		}
	}
	toMap["index"] = index

	return nil
}

func (film *filmYaml) MarshalYAML() (interface{}, error) {
	toMap := make(map[string]interface{})
	if err := film.PushToMap(toMap); err != nil {
		return nil, fmt.Errorf("pushing film to map: %w", err)
	}

	return toMap, nil
}

func (film *filmYaml) UnmarshalYAML(value *yaml.Node) error {
	if value.Kind != yaml.MappingNode {
		return fmt.Errorf("not a mapping node for film")
	}

	// Simpler to code than pulling everything bit by bit from the value object.
	// The latter might be faster, however.
	temp := make(map[string]interface{})
	if err := value.Decode(temp); err != nil {
		return fmt.Errorf("decoding film to temp: %w", err)
	}

	if err := film.PullFromMap(temp); err != nil {
		return fmt.Errorf("pulling film from map: %w", err)
	}

	return nil
}

//////////////////////////////////////////////////////////////////////////

func (suite *YamlTestSuite) TestMarshalCycle() {
	suite.film.Open(suite.mapper)
	defer suite.film.Close()

	bytes, err := yaml.Marshal(suite.film)
	suite.Assert().NoError(err)

	fmt.Print("--- TestMarshalCycle ---------------\n", string(bytes), "------------------------------------\n")

	var film filmYaml
	film.Open(suite.mapper)
	suite.Assert().NoError(yaml.Unmarshal(bytes, &film))
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

func (suite *YamlTestSuite) TestLoadFromString() {
	base, _ := NewConverter(suite.mapper)
	loaded, err := base.LoadFromString(simpleYaml)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(loaded)
	film, ok := loaded.(*filmYaml)
	suite.Assert().True(ok)
	suite.Assert().NotNil(film)
	suite.Assert().Equal("Blockbuster Movie", film.Name)
	suite.checkAlpha(film.Lead)
	suite.Assert().NotNil(film.Cast)
	suite.Assert().Len(film.Cast, 2)
	suite.checkAlpha(film.Cast[0])
	suite.checkBravo(film.Cast[1])
	suite.Assert().NotNil(film.Index)
	suite.Assert().Len(film.Index, 2)
	suite.checkAlpha(film.Index["Lucky, Lance"])
	suite.checkBravo(film.Index["Queue, Susie"])
}

func (suite *YamlTestSuite) TestMarshalFileCycle() {
	marshal, _ := NewConverter(suite.mapper)
	file, err := ioutil.TempFile("", "*.test.yaml")
	suite.Assert().NoError(err)
	suite.Assert().NotNil(file)
	fileName := file.Name()
	// Go ahead and close it, just needed the file name.
	suite.Assert().NoError(file.Close())

	film := suite.makeTestFilm()
	suite.Assert().NoError(marshal.SaveToFile(film, fileName))

	reloaded, err := marshal.LoadFromFile(fileName)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(reloaded)
	suite.Assert().Equal(film, reloaded)
}

func (suite *YamlTestSuite) TestMarshalStringCycle() {
	marshal, _ := NewConverter(suite.mapper)
	film := suite.makeTestFilm()
	str, err := marshal.SaveToString(film)
	suite.Assert().NoError(err)
	suite.NotZero(str)

	fmt.Print("--- TestMarshalStringCycle ---------\n", str, "------------------------------------\n")

	reloaded, err := marshal.LoadFromString(str)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(reloaded)
	suite.Assert().Equal(film, reloaded)
}

//////////////////////////////////////////////////////////////////////////

const simpleYaml = `
<type>: '[testYAML]filmYaml'
name:   'Blockbuster Movie'
lead: {
  <type>: '[test]Alpha',
  name: 'Lance Lucky',
  percent: 23.79,
  extra: 'Yaaaa!'
}
cast:
- {
    <type>: '[test]Alpha',
    name: 'Lance Lucky',
    percent: 23.79,
    extra: false
  }
- {
    <type>: '[test]Bravo',
    finished: true,
    iterations: 13,
    extra: 'gibbering ghostwhistle'
  }
index: {
  'Lucky, Lance': {
    <type>: '[test]Alpha',
    name: 'Lance Lucky',
    percent: 23.79,
    extra: 'marshmallow stars'
  },
  'Queue, Susie': {
    <type>: '[test]Bravo',
    finished: true,
    iterations: 13,
    extra: 19.57
  }
}
`

func (suite *YamlTestSuite) checkAlpha(act test.Actor) {
	suite.Assert().NotNil(act)
	alf, ok := act.(*test.Alpha)
	suite.Assert().True(ok)
	suite.Assert().NotNil(alf)
	suite.Assert().Equal("Lance Lucky", alf.Name)
	suite.Assert().Equal(float32(23.79), alf.Percent)
	suite.Assert().Empty(alf.Extra())
}

func (suite *YamlTestSuite) checkBravo(act test.Actor) {
	suite.Assert().NotNil(act)
	bra, ok := act.(*test.Bravo)
	suite.Assert().True(ok)
	suite.Assert().NotNil(bra)
	suite.Assert().True(bra.Finished)
	suite.Assert().Equal(13, bra.Iterations)
	suite.Assert().Empty(bra.Extra())
}

func (suite *YamlTestSuite) makeTestFilm() *filmYaml {
	actor1 := &test.Alpha{
		Name:    "Goober Snoofus",
		Percent: 13.23,
	}
	actor2 := &test.Bravo{
		Finished:   true,
		Iterations: 1957,
	}
	return &filmYaml{
		Name: "",
		Lead: actor1,
		Cast: []test.Actor{
			actor1,
			actor2,
		},
		Index: map[string]test.Actor{
			"Snoofus, Goober": actor1,
			"Snarly, Booger":  actor2,
		},
	}
}
