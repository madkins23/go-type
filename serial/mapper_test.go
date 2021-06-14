package serial

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/madkins23/go-type/datamap"
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

func (tp *TestProduct) Marshal() (map[string]interface{}, error) {
	return datamap.Marshal(tp)
}

func (tp *TestProduct) Unmarshal(fromMap map[string]interface{}) error {
	// TODO: Not sure how to Unmarshal _into_ tp?!?
	return datamap.Unmarshal(fromMap, tp)
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
	suite.mapper = NewMapper(suite.registry)
}

func TestMapperSuite(t *testing.T) {
	suite.Run(t, new(mapperTestSuite))
}

//////////////////////////////////////////////////////////////////////////

func (suite *mapperTestSuite) TestMarshal() {
	suite.Assert().NoError(suite.registry.Register(&TestProduct{}))
	m, err := suite.mapper.Marshal(&TestProduct{Name: prodName, Price: price, extra: "nothing to see here"})
	suite.Assert().NoError(err)
	suite.Assert().NotNil(m)
	fmt.Printf("toMap: %#v\n", m)
	suite.Assert().Len(m, 3)
	suite.Assert().Equal(suite.packagePath+"/"+structName, m[reg.TypeField])
	suite.Assert().Equal(prodName, m["Name"])
	suite.Assert().Equal(price, m["Price"])
}

func (suite *mapperTestSuite) TestUnmarshal() {
	suite.Assert().NoError(suite.registry.Register(&TestProduct{}))
	example, err := suite.mapper.Unmarshal(map[string]interface{}{
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
