package reg

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/madkins23/go-type/test"
)

//////////////////////////////////////////////////////////////////////////

func ExampleRegistry_Alias() {
	registry := NewRegistry()
	if registry.Alias("[alpha]", &test.Alpha{}) == nil {
		fmt.Println("Aliased")
	}
	// output: Aliased
}

func ExampleRegistry_Register() {
	registry := NewRegistry()
	if registry.Register(&test.Alpha{}) == nil {
		fmt.Println("Registered")
	}
	// output: Registered
}

//////////////////////////////////////////////////////////////////////////

type registryTestSuite struct {
	suite.Suite
	registry Registry
	reg      *registry
}

func (suite *registryTestSuite) SetupTest() {
	// See json_test.go for these definitions.
	test.CopyMapFromItemFn = test.CopyMapFromItemJSON
	test.CopyItemFromMapFn = test.CopyItemFromMapJSON

	suite.registry = NewRegistry()
	var ok bool
	suite.reg, ok = suite.registry.(*registry)
	suite.Assert().True(ok)
}

func (suite *registryTestSuite) TearDownSuite() {
	test.CopyMapFromItemFn = nil
	test.CopyItemFromMapFn = nil
}

func TestRegistrySuite(t *testing.T) {
	suite.Run(t, new(registryTestSuite))
}

//////////////////////////////////////////////////////////////////////////

func (suite *registryTestSuite) TestNewRegistry() {
	suite.Assert().NotNil(suite.registry)
	suite.Assert().NotNil(suite.reg.byName)
	suite.Assert().NotNil(suite.reg.byType)
}

func (suite *registryTestSuite) TestAlias() {
	example := &test.Alpha{}
	err := suite.registry.Alias("badPackage", &example)
	suite.Assert().Error(err)
	suite.Assert().Contains(err.Error(), "no package path")
	suite.Assert().Empty(suite.reg.aliases)
	err = suite.registry.Alias("x", example)
	suite.Assert().NoError(err)
	suite.Assert().Len(suite.reg.aliases, 1)
	err = suite.registry.Alias("x", example)
	suite.Assert().Error(err)
	suite.Assert().Contains(err.Error(), "can't redefine alias")
}

func (suite *registryTestSuite) TestRegister() {
	example := &test.Alpha{}
	err := suite.registry.Register(&example)
	suite.Assert().Error(err)
	suite.Assert().Contains(err.Error(), "no path for type")
	suite.Assert().Empty(suite.reg.byName)
	suite.Assert().Empty(suite.reg.byType)
	err = suite.registry.Register(example)
	suite.Assert().NoError(err)
	suite.Assert().Len(suite.reg.byName, 1)
	suite.Assert().Len(suite.reg.byType, 1)
	err = suite.registry.Register(example)
	suite.Assert().Error(err)
	suite.Assert().Contains(err.Error(), "previous registration")
}

func (suite *registryTestSuite) TestNameFor() {
	example := &test.Alpha{}
	suite.Assert().NoError(suite.registry.Register(example))
	exType, err := suite.registry.NameFor(example)
	suite.Assert().NoError(err)
	suite.Assert().Equal(test.PackageName+"/Alpha", exType)
}

func (suite *registryTestSuite) TestMake() {
	example := &test.Alpha{}
	suite.Assert().NoError(suite.registry.Register(example))
	item, err := suite.registry.Make(test.PackageName + "/Alpha")
	suite.Assert().NoError(err)
	suite.Assert().NotNil(item)
	suite.Assert().IsType(example, item)
}

func (suite *registryTestSuite) TestConverItemToMap() {
	suite.Assert().NoError(suite.registry.Register(&test.Alpha{}))
	m, err := suite.registry.ConvertItemToMap(test.NewAlpha("Goober Snoofus", 17.23, "nothing to see here"))
	suite.Assert().NoError(err)
	suite.Assert().NotNil(m)
	suite.Assert().Len(m, 3)
	suite.Assert().Equal(test.PackageName+"/Alpha", m[TypeField])
	suite.Assert().Equal("Goober Snoofus", m["Name"])
	suite.Assert().Equal(17.23, m["Percent"])
}

func (suite *registryTestSuite) TestCreateItemFromMap() {
	suite.Assert().NoError(suite.registry.Register(&test.Alpha{}))
	example, err := suite.registry.CreateItemFromMap(map[string]interface{}{
		TypeField: test.PackageName + "/Alpha",
		"Name":    "Goober Snoofus",
		"Percent": 17.23,
		"extra":   "nothing to see here",
	})
	suite.Assert().NoError(err)
	suite.Assert().NotNil(example)
	suite.Assert().IsType(&test.Alpha{}, example)
	suite.Assert().Equal(&test.Alpha{
		Name:    "Goober Snoofus",
		Percent: 17.23,
	}, example)
}

func (suite *registryTestSuite) TestCycleSimple() {
	example := &test.Alpha{}
	suite.Assert().NoError(suite.registry.Register(example))
	registration := suite.reg.byType[reflect.TypeOf(example).Elem()]
	suite.Assert().NotNil(registration)
	suite.Assert().Len(registration.allNames, 1)
	name, err := suite.registry.NameFor(example)
	suite.Assert().NoError(err)
	suite.Assert().Equal(test.PackageName+"/Alpha", name)
	object, err := suite.registry.Make(name)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(object)
	suite.Assert().Equal(reflect.TypeOf(example), reflect.TypeOf(object))
}

func (suite *registryTestSuite) TestCycleAlias() {
	example := &test.Alpha{}
	suite.Assert().NoError(suite.registry.Alias("typeUtils", example))
	suite.Assert().NoError(suite.registry.Register(example))
	exType := reflect.TypeOf(example)
	if exType.Kind() == reflect.Ptr {
		exType = exType.Elem()
	}
	registration, ok := suite.reg.byType[exType]
	suite.Assert().True(ok)
	suite.Assert().NotNil(registration)

	suite.Assert().Len(registration.allNames, 2)
	name, err := suite.registry.NameFor(example)
	suite.Assert().NoError(err)
	suite.Assert().Equal("[typeUtils]Alpha", name)
	object, err := suite.registry.Make(name)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(object)
	suite.Assert().Equal(reflect.TypeOf(example), reflect.TypeOf(object))
}

func (suite *registryTestSuite) TestGenNames() {
	example := &test.Alpha{}
	name, aliases, err := suite.reg.GenNames(example, false)
	suite.Assert().NoError(err)
	suite.Assert().Equal(test.PackageName+"/Alpha", name)
	suite.Assert().Nil(aliases)
	name, aliases, err = suite.reg.GenNames(example, true)
	suite.Assert().NoError(err)
	suite.Assert().Equal(test.PackageName+"/Alpha", name)
	suite.Assert().NotNil(aliases)
	suite.Assert().Empty(aliases)

	suite.Assert().NoError(suite.registry.Alias("typeUtils", example))
	name, aliases, err = suite.reg.GenNames(example, false)
	suite.Assert().NoError(err)
	suite.Assert().Equal(test.PackageName+"/Alpha", name)
	suite.Assert().Nil(aliases)
	name, aliases, err = suite.reg.GenNames(example, true)
	suite.Assert().NoError(err)
	suite.Assert().Equal("[typeUtils]Alpha", name)
	suite.Assert().NotNil(aliases)
	suite.Assert().Len(aliases, 1)

	name, aliases, err = suite.reg.GenNames(&example, true)
	suite.Assert().Error(err)
	suite.Assert().Contains(err.Error(), "no path for type")
	name, aliases, err = suite.reg.GenNames(1, true)
	suite.Assert().Error(err)
	suite.Assert().Contains(err.Error(), "no path for type")
}

func (suite *registryTestSuite) TestGenTypeName() {
	example := &test.Alpha{}
	name, err := genNameFromInterface(example)
	suite.Assert().NoError(err)
	suite.Assert().Equal(test.PackageName+"/Alpha", name)

	_, err = genNameFromInterface(&example)
	suite.Assert().Error(err)
	suite.Assert().Contains(err.Error(), "no path for type")
	_, err = genNameFromInterface(1)
	suite.Assert().Error(err)
	suite.Assert().Contains(err.Error(), "no path for type")
}

//////////////////////////////////////////////////////////////////////////

// Make sure io.ReadSeeker works the way we think.
func (suite *registryTestSuite) TestReadSeeker() {
	var err error
	var readSeek io.ReadSeeker

	stringReader := strings.NewReader(jabberwocky)
	suite.Assert().NotNil(stringReader)
	readSeek = stringReader
	suite.Assert().NotNil(readSeek)

	file, err := ioutil.TempFile("", "*.test")
	defer func() {
		suite.Assert().NoError(file.Close())
		suite.Assert().NoError(os.Remove(file.Name()))
	}()
	suite.Assert().NoError(err)
	suite.Assert().NotNil(file)
	readSeek = file
	suite.Assert().NotNil(readSeek)

	_, err = file.Write([]byte(jabberwocky))
	suite.Assert().NoError(err)
	where, err := file.Seek(0, io.SeekStart)
	suite.Assert().NoError(err)
	suite.Assert().Zero(where)
	bytes, err := ioutil.ReadAll(file)
	suite.Assert().NoError(err)
	suite.Assert().Equal(jabberwocky, string(bytes))
}

//////////////////////////////////////////////////////////////////////////

const jabberwocky = `
'Twas brillig, and the slithy toves.
Did gyre and gimble in the wabe:
All mimsy were the borogoves,
And the mome raths outgrabe.
`
