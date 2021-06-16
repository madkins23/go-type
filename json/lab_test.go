package json

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/madkins23/go-type/reg"
	"github.com/madkins23/go-type/test"
	"github.com/stretchr/testify/suite"
)

// Trial and error approach to simplifying things.

var showAccount = true

type LabTestSuite struct {
	suite.Suite
}

func (suite *LabTestSuite) SetupSuite() {
	suite.Assert().NoError(reg.AddAlias("test", test.Account{}), "creating test alias")
	suite.Assert().NoError(reg.Register(&test.Stock{}))
	suite.Assert().NoError(reg.Register(&test.Bond{}))
	suite.Assert().NoError(reg.AddAlias("labTest", Account{}), "creating test alias")
	suite.Assert().NoError(reg.Register(&Account{}))
}

func (suite *LabTestSuite) SetupTest() {
}

func TestLabSuite(t *testing.T) {
	suite.Run(t, new(LabTestSuite))
}

//////////////////////////////////////////////////////////////////////////

// TestMarshalCycle verifies the JSON Marshal/Unmarshal works as expected.
func (suite *LabTestSuite) TestMarshalCycle() {
	account := MakeAccount()

	// TODO: Still HTML-escaping, damn it.
	marshaled, err := json.Marshal(account)
	suite.Assert().NoError(err)

	if showAccount {
		var buf bytes.Buffer
		suite.NoError(json.Indent(&buf, marshaled, "", "  "))
		fmt.Println(buf.String())
	}

	var newAccount test.Account
	suite.Assert().NoError(json.Unmarshal(marshaled, &newAccount))

	//suite.Assert().NotEqual(suite.film, &film) // fails due to unexported field 'extra'
	//for _, act := range suite.film.Cast {
	//	// Remove unexported field.
	//	if a, ok := act.(*test.Alpha); ok {
	//		a.ClearExtra()
	//	} else if b, ok := act.(*test.Bravo); ok {
	//		b.ClearExtra()
	//	}
	//}
	//suite.Assert().Equal(suite.film, &film) // succeeds now that unexported fields are gone.
}

//////////////////////////////////////////////////////////////////////////

func MakeAccount() *Account {
	account := &Account{}
	account.MakeFake()
	return account
}

type Account struct {
	test.Account
}

func (a *Account) MarshalJSON() ([]byte, error) {
	work := &strings.Builder{}
	if _, err := fmt.Fprintf(work, "{"); err != nil {
		return nil, fmt.Errorf("write open brace: %w", err)
	}

	if a.Favorite != nil {
		if err := MarshalFieldItem("favorite", a.Favorite, false, work); err != nil {
			return nil, fmt.Errorf("marshal field item: %w", err)
		}
	}

	if _, err := fmt.Fprintf(work, "}"); err != nil {
		return nil, fmt.Errorf("write close brace: %w", err)
	}

	return []byte(work.String()), nil

}

//////////////////////////////////////////////////////////////////////////

var simpleLabJson = `
{
  "Favorite": {
    "Market": "NASDAQ",
    "Symbol": "COST",
    "Name": "Costco",
    "Position": 10,
    "Value": 400
  },
  "Positions": [
    {
      "Market": "NASDAQ",
      "Symbol": "COST",
      "Name": "Costco",
      "Position": 10,
      "Value": 400
    },
    {
      "Market": "NYSE",
      "Symbol": "WMT",
      "Name": "Walmart",
      "Position": 20,
      "Value": 150
    },
    {
      "Source": "Treasury",
      "Name": "T-Bill",
      "Value": 1000,
      "Interest": 0.75,
      "Duration": 31536000000000000,
      "Expires": "2021-07-08T11:34:13.556165402-07:00"
    }
  ],
  "Lookup": {
    "COST": {
      "Market": "NASDAQ",
      "Symbol": "COST",
      "Name": "Costco",
      "Position": 10,
      "Value": 400
    },
    "WMT": {
      "Market": "NYSE",
      "Symbol": "WMT",
      "Name": "Walmart",
      "Position": 20,
      "Value": 150
    }
  }
}
`
