package convert

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type Person struct {
	Name   string
	Age    int
	Emails []string
	Extra  map[string]string
}

func (p *Person) PushToMap(toMap map[string]interface{}) error {
	return PushItemToMap(p, toMap)
}

func (p *Person) PullFromMap(fromMap map[string]interface{}) error {
	return PullItemFromMap(p, fromMap)
}

const age = 23
const name = "Goober Snoofus"

func emails() []string {
	return []string{"goober@snoofus.nom", "gSnoofus@evilcorp.com", "gSnoof@example.com"}
}

func extra() map[string]string {
	return map[string]string{
		"twitter": "mitchellh",
	}
}

var mapped map[string]interface{}
var structured Person

func init() {
	mapped = map[string]interface{}{
		"Name":   name,
		"Age":    age,
		"Emails": emails(),
		"Extra":  extra(),
	}
	structured = Person{
		Name:   name,
		Age:    age,
		Emails: emails(),
		Extra:  extra(),
	}
}

func TestMapperBase_PullFromMap(t *testing.T) {
	var result Person
	err := result.PullFromMap(mapped)
	assert.NoError(t, err)
	//fmt.Printf("Pulled: %#v\n  from: %#v\n", result, mapped)
	assert.Equal(t, structured, result)
}

func TestMapperBase_PushToMap(t *testing.T) {
	result := make(map[string]interface{})
	err := structured.PushToMap(result)
	assert.NoError(t, err)
	//fmt.Printf("Pushed: %#v\n    to: %#v\n", structured, result)
	assert.Equal(t, mapped, result)
}
