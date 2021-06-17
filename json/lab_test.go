package json

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/suite"

	"github.com/madkins23/go-type/data"
	"github.com/madkins23/go-type/reg"
	"github.com/madkins23/go-type/test"
)

// Trial and error approach to simplifying things.

var showAccount = false

type LabTestSuite struct {
	suite.Suite
}

func (suite *LabTestSuite) SetupSuite() {
	suite.Require().NoError(reg.AddAlias("test", test.Account{}), "creating test alias")
	suite.Require().NoError(reg.Register(&test.Stock{}))
	suite.Require().NoError(reg.Register(&test.Bond{}))
	suite.Require().NoError(reg.AddAlias("labTest", Account{}), "creating test alias")
	suite.Require().NoError(reg.Register(&Account{}))
}

func (suite *LabTestSuite) SetupTest() {
}

func TestLabSuite(t *testing.T) {
	suite.Run(t, new(LabTestSuite))
}

// TestMarshalCycle verifies the JSON Marshal/Unmarshal works as expected.
func (suite *LabTestSuite) TestMarshalCycle() {
	account := MakeAccount()

	marshaled, err := json.Marshal(account)
	suite.Require().NoError(err)
	if showAccount {
		var buf bytes.Buffer
		suite.Require().NoError(json.Indent(&buf, marshaled, "", "  "))
		fmt.Println(buf.String())
	}
	suite.Assert().Contains(string(marshaled), "TypeName")
	suite.Assert().Contains(string(marshaled), "[test]Stock")
	suite.Assert().Contains(string(marshaled), "[test]Bond")

	var newAccount Account
	suite.Require().NoError(json.Unmarshal(marshaled, &newAccount))
	if showAccount {
		fmt.Println("---------------------------")
		spew.Dump(newAccount)
	}

	suite.Assert().NotEqual(account, newAccount)
	account.Favorite.ClearPrivateFields()
	for _, position := range account.Positions {
		position.ClearPrivateFields()
	}
	for _, position := range account.Lookup {
		position.ClearPrivateFields()
	}
	// Succeeds now that unexported (private) fields are gone.
	suite.Assert().Equal(account, &newAccount)
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
	temp, err := data.Marshal(a)
	if err != nil {
		return nil, fmt.Errorf("create data.Map: %w", err)
	}

	// Wrap objects referenced by interface fields.
	if compAccount, found := temp["Account"]; !found {
		return nil, fmt.Errorf("no composited Account")
	} else if labAccount, ok := compAccount.(map[string]interface{}); !ok {
		return nil, fmt.Errorf("bad composited Account:  %#v", labAccount)
	} else {
		if a.Favorite != nil {
			if labAccount["Favorite"], err = WrapItem(a.Favorite); err != nil {
				return nil, fmt.Errorf("wrap favorite: %w", err)
			}
		}
		if a.Positions != nil {
			fixed := make([]*Wrapper, len(a.Positions))
			for i, pos := range a.Positions {
				if fixed[i], err = WrapItem(pos); err != nil {
					return nil, fmt.Errorf("wrap Positions item: %w", err)
				}
			}
			labAccount["Positions"] = fixed
		}
		if a.Lookup != nil {
			fixed := make(map[string]*Wrapper, len(a.Lookup))
			for k, pos := range a.Lookup {
				if fixed[k], err = WrapItem(pos); err != nil {
					return nil, fmt.Errorf("wrap Lookup item: %w", err)
				}
			}
			labAccount["Lookup"] = fixed
		}
	}

	build := &strings.Builder{}
	encoder := json.NewEncoder(build)
	encoder.SetEscapeHTML(false)
	if err = encoder.Encode(temp); err != nil {
		return nil, fmt.Errorf("encode wrapped item to JSON: %w", err)
	}
	return []byte(build.String()), nil
}

type loadAccount struct {
	Account struct {
		Favorite  *Wrapper
		Positions []*Wrapper
		Lookup    map[string]*Wrapper
	}
}

func (a *Account) getInvestment(w *Wrapper) (test.Investment, error) {
	var ok bool
	var investment test.Investment
	if w != nil {
		if item, err := w.Unwrap(); err != nil {
			return nil, fmt.Errorf("unwrap item: %w", err)
		} else if investment, ok = item.(test.Investment); !ok {
			return nil, fmt.Errorf("item %#v not Investment", item)
		}
	}

	return investment, nil
}

func (a *Account) UnmarshalJSON(marshaled []byte) error {
	loader := &loadAccount{}
	if err := json.Unmarshal(marshaled, loader); err != nil {
		return fmt.Errorf("unmarshal to loader: %w", err)
	}

	var err error
	if a.Favorite, err = a.getInvestment(loader.Account.Favorite); err != nil {
		return fmt.Errorf("get Investment from Favorite")
	}
	if loader.Account.Positions != nil {
		fixed := make([]test.Investment, len(loader.Account.Positions))
		for i, wPos := range loader.Account.Positions {
			if fixed[i], err = a.getInvestment(wPos); err != nil {
				return fmt.Errorf("get Investment from Positions")
			}
		}
		a.Positions = fixed
	}
	if loader.Account.Lookup != nil {
		fixed := make(map[string]test.Investment, len(loader.Account.Lookup))
		for key, wPos := range loader.Account.Lookup {
			if fixed[key], err = a.getInvestment(wPos); err != nil {
				return fmt.Errorf("get Investment from Lookup")
			}
		}
		a.Lookup = fixed
	}

	return nil
}
