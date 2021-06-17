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

	// TODO: use a simplistic account where the positions are all Wrappers?

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

func (a *Account) UnmarshalJSON(marshaled []byte) error {
	dataMap := make(data.Map)
	decoder := json.NewDecoder(strings.NewReader(string(marshaled)))
	if err := decoder.Decode(&dataMap); err != nil {
		return fmt.Errorf("decode JSON: %w", err)
	}

	if compAccount, found := dataMap["Account"]; !found {
		return fmt.Errorf("no composited Account")
	} else if labAccount, ok := compAccount.(map[string]interface{}); !ok {
		return fmt.Errorf("bad composited Account:  %#v", labAccount)
	} else {
		var err error
		if labAccount["Favorite"] != nil {
			if favorite, ok := labAccount["Favorite"].(map[string]interface{}); !ok {
				return fmt.Errorf("favorite not data map")
			} else if labAccount["Favorite"], err = UnwrapItem(favorite); err != nil {
				return fmt.Errorf("unwrap item: %w", err)
			}
		}
		if labAccount["Positions"] != nil {
			if positions, ok := (labAccount["Positions"]).([]interface{}); !ok {
				return fmt.Errorf("positions not array")
			} else {
				fixed := make([]interface{}, len(positions))
				for i, pos := range positions {
					if position, ok := pos.(map[string]interface{}); !ok {
						return fmt.Errorf("position not data map")
					} else if fixed[i], err = UnwrapItem(position); err != nil {
						return fmt.Errorf("unwrap item: %w", err)
					}
				}
				labAccount["Positions"] = fixed
			}
		}
		if labAccount["Lookup"] != nil {
			if lookup, ok := (labAccount["Lookup"]).(map[string]interface{}); !ok {
				return fmt.Errorf("lookup not map")
			} else {
				fixed := make(map[string]interface{}, len(lookup))
				for key, pos := range lookup {
					if position, ok := pos.(map[string]interface{}); !ok {
						return fmt.Errorf("lookup position not data map")
					} else if fixed[key], err = UnwrapItem(position); err != nil {
						return fmt.Errorf("unwrap item: %w", err)
					}
				}
				labAccount["Lookup"] = fixed
			}
		}
		if err := data.Unmarshal(dataMap, a); err != nil {
			return fmt.Errorf("unmarshal map data: %w", err)
		}
	}

	return nil
}
