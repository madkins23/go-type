package serial

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/madkins23/go-type/convert"

	"github.com/madkins23/go-type/reg"
	"github.com/stretchr/testify/suite"
)

//////////////////////////////////////////////////////////////////////////

const prodName = "Organic Vulture Giblets"
const price = float32(17.23)

const structName = "TestProduct"

type TestProduct struct {
	Name  string
	Price float32
	extra string
}

func (tp *TestProduct) PushToMap(toMap map[string]interface{}) error {
	return convert.PushItemToMap(tp, toMap)
}

func (tp *TestProduct) PullFromMap(fromMap map[string]interface{}) error {
	return convert.PullItemFromMap(tp, fromMap)
}

//////////////////////////////////////////////////////////////////////////

type mapperTestSuite struct {
	suite.Suite
	mapper      Mapper
	registry    reg.Registry
	packagePath string
}

func (suite *mapperTestSuite) SetupSuite() {
	exampleType := reflect.TypeOf(TestProduct{})
	suite.Assert().NotNil(exampleType)
	suite.packagePath = exampleType.PkgPath()
	suite.Assert().NotEmpty(suite.packagePath)
}

func (suite *mapperTestSuite) SetupTest() {
	suite.registry = reg.NewRegistry()
	var err error
	suite.mapper, err = NewMapper(suite.registry)
	suite.Assert().NoError(err)
}

func TestMapperSuite(t *testing.T) {
	suite.Run(t, new(mapperTestSuite))
}

//////////////////////////////////////////////////////////////////////////

func (suite *mapperTestSuite) TestConvertItemToMap() {
	suite.Assert().NoError(suite.registry.Register(&TestProduct{}))
	m, err := suite.mapper.ConvertItemToMap(&TestProduct{Name: prodName, Price: price, extra: "nothing to see here"})
	suite.Assert().NoError(err)
	suite.Assert().NotNil(m)
	fmt.Printf("toMap: %#v\n", m)
	suite.Assert().Len(m, 3)
	suite.Assert().Equal(suite.packagePath+"/"+structName, m[reg.TypeField])
	suite.Assert().Equal(prodName, m["Name"])
	suite.Assert().Equal(price, m["Price"])
}

func (suite *mapperTestSuite) TestCreateItemFromMap() {
	suite.Assert().NoError(suite.registry.Register(&TestProduct{}))
	example, err := suite.mapper.CreateItemFromMap(map[string]interface{}{
		reg.TypeField: suite.packagePath + "/" + structName,
		"Name":        prodName,
		"Price":       price,
		"extra":       "nothing to see here",
	})
	suite.Assert().NoError(err)
	suite.Assert().NotNil(example)
	suite.Assert().IsType(&TestProduct{}, example)
	suite.Assert().Equal(&TestProduct{
		Name:  prodName,
		Price: price,
	}, example)
}
