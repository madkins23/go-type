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
)

//////////////////////////////////////////////////////////////////////////

func ExampleRegistry_AddAlias() {
	registry := NewRegistry()
	if registry.AddAlias("[alpha]", &Alpha{}) == nil {
		fmt.Println("Aliased")
	}
	// output: Aliased
}

func ExampleRegistry_Register() {
	registry := NewRegistry()
	if registry.Register(&Alpha{}) == nil {
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
	suite.registry = NewRegistry()
	var ok bool
	suite.reg, ok = suite.registry.(*registry)
	suite.Assert().True(ok)
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
	example := &Alpha{}
	err := suite.registry.AddAlias("badPackage", &example)
	suite.Assert().Error(err)
	suite.Assert().Contains(err.Error(), "no package path")
	suite.Assert().Empty(suite.reg.aliases)
	err = suite.registry.AddAlias("x", example)
	suite.Assert().NoError(err)
	suite.Assert().Len(suite.reg.aliases, 1)
	err = suite.registry.AddAlias("x", example)
	suite.Assert().Error(err)
	suite.Assert().Contains(err.Error(), "can't redefine alias")
}

func (suite *registryTestSuite) TestRegister() {
	example := &Alpha{}
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
	type localStruct struct {
		something string
	}
	err = suite.registry.Register(&localStruct{})
	suite.Assert().Error(err)
	suite.Assert().Contains(err.Error(), "is private")
}

func (suite *registryTestSuite) TestNameFor() {
	example := &Alpha{}
	suite.Assert().NoError(suite.registry.Register(example))
	exType, err := suite.registry.NameFor(example)
	suite.Assert().NoError(err)
	suite.Assert().Equal(packageName+"/Alpha", exType)
}

func (suite *registryTestSuite) TestNameForNilItem() {
	exType, err := suite.registry.NameFor(nil)
	suite.Assert().Equal("", exType)
	suite.Assert().ErrorIs(err, errItemIsNil)
}

func (suite *registryTestSuite) TestMake() {
	example := &Alpha{}
	suite.Assert().NoError(suite.registry.Register(example))
	item, err := suite.registry.Make(packageName + "/Alpha")
	suite.Assert().NoError(err)
	suite.Assert().NotNil(item)
	suite.Assert().IsType(example, item)
}

func (suite *registryTestSuite) TestCycleSimple() {
	example := &Alpha{}
	suite.Assert().NoError(suite.registry.Register(example))
	registration := suite.reg.byType[reflect.TypeOf(example).Elem()]
	suite.Assert().NotNil(registration)
	suite.Assert().Len(registration.allNames, 1)
	name, err := suite.registry.NameFor(example)
	suite.Assert().NoError(err)
	suite.Assert().Equal(packageName+"/Alpha", name)
	object, err := suite.registry.Make(name)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(object)
	suite.Assert().Equal(reflect.TypeOf(example), reflect.TypeOf(object))
}

func (suite *registryTestSuite) TestCycleAlias() {
	example := &Alpha{}
	suite.Assert().NoError(suite.registry.AddAlias("typeUtils", example))
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
	example := &Alpha{}
	name, aliases, err := suite.reg.genNames(example, false)
	suite.Assert().NoError(err)
	suite.Assert().Equal(packageName+"/Alpha", name)
	suite.Assert().Nil(aliases)
	name, aliases, err = suite.reg.genNames(example, true)
	suite.Assert().NoError(err)
	suite.Assert().Equal(packageName+"/Alpha", name)
	suite.Assert().NotNil(aliases)
	suite.Assert().Empty(aliases)

	suite.Assert().NoError(suite.registry.AddAlias("typeUtils", example))
	name, aliases, err = suite.reg.genNames(example, false)
	suite.Assert().NoError(err)
	suite.Assert().Equal(packageName+"/Alpha", name)
	suite.Assert().Nil(aliases)
	name, aliases, err = suite.reg.genNames(example, true)
	suite.Assert().NoError(err)
	suite.Assert().Equal("[typeUtils]Alpha", name)
	suite.Assert().NotNil(aliases)
	suite.Assert().Len(aliases, 1)

	_, _, err = suite.reg.genNames(&example, true)
	suite.Assert().Error(err)
	suite.Assert().Contains(err.Error(), "no path for type")
	_, _, err = suite.reg.genNames(1, true)
	suite.Assert().Error(err)
	suite.Assert().Contains(err.Error(), "no path for type")
}

func (suite *registryTestSuite) TestGenTypeName() {
	example := &Alpha{}
	name, err := genNameFromInterface(example)
	suite.Assert().NoError(err)
	suite.Assert().Equal(packageName+"/Alpha", name)

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
