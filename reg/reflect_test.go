package reg

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/suite"
)

// These tests confirm the developer's understanding of how Go works.
// More specifically how the Go reflection mechanism works.

var (
	a        = Alpha{Name: "Hubert", Number: 17.23}
	b        = Bravo{Finished: true, Iterations: 79}
	c        = a
	ai Stuff = &a
	bi Stuff = &b
	ci Stuff = &c
)

//////////////////////////////////////////////////////////////////////////

type ReflectTestSuite struct {
	suite.Suite
}

func (suite *ReflectTestSuite) SetupTest() {
}
func TestReflectSuite(t *testing.T) {
	suite.Run(t, new(ReflectTestSuite))
}

//////////////////////////////////////////////////////////////////////////
// Verify method for determining path of package via an object defined therein.

func (suite *ReflectTestSuite) TestPackagePath() {
	suite.Assert().Equal(reflect.TypeOf(Alpha{}).PkgPath(), packageName)
}

//////////////////////////////////////////////////////////////////////////
// MakeFake certain reflect.Type supports equivalence testing and use as map key.
// Note that this does NOT work for types.Type, which is a different thing.

func (suite *ReflectTestSuite) TestTypeEquivalence() {
	suite.Assert().NotEqual(a, b)
	suite.Assert().NotEqual(b, c)
	suite.Assert().Equal(a, c)
}

func (suite *ReflectTestSuite) TestInterfaceEquivalence() {
	suite.Assert().NotEqual(ai, bi)
	suite.Assert().NotEqual(bi, ci)
	suite.Assert().Equal(ai, ci)
}

func (suite *ReflectTestSuite) TestMap() {
	lookup := make(map[reflect.Type]string)
	lookup[reflect.TypeOf(a)] = "alpha"
	lookup[reflect.TypeOf(b)] = "bravo"
	lookup[reflect.TypeOf(c)] = "charlie" // overrides "alpha" since a == c
	suite.Assert().Equal("charlie", lookup[reflect.TypeOf(a)])
	suite.Assert().Equal("bravo", lookup[reflect.TypeOf(b)])
	suite.Assert().Equal("charlie", lookup[reflect.TypeOf(c)])
}

func (suite *ReflectTestSuite) TestMapInterface() {
	lookup := make(map[reflect.Type]string)
	lookup[reflect.TypeOf(ai)] = "alpha"
	lookup[reflect.TypeOf(bi)] = "bravo"
	lookup[reflect.TypeOf(ci)] = "charlie" // overrides "alpha" since ai == ci
	suite.Assert().Equal("charlie", lookup[reflect.TypeOf(ai)])
	suite.Assert().Equal("bravo", lookup[reflect.TypeOf(bi)])
	suite.Assert().Equal("charlie", lookup[reflect.TypeOf(ci)])
}

//////////////////////////////////////////////////////////////////////////

const packageName = "github.com/madkins23/go-type/reg"

type Stuff interface {
	Info() string
}

type Alpha struct {
	Name   string
	Number float32
}

func (a *Alpha) Info() string {
	return fmt.Sprintf("%s: %f", a.Name, a.Number)
}

type Bravo struct {
	Finished   bool
	Iterations int
}

func (b *Bravo) Info() string {
	return fmt.Sprintf("%t: %d", b.Finished, b.Iterations)
}
