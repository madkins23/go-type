package datamap

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

func (p *Person) Marshal() (map[string]interface{}, error) {
	return Marshal(p)
}

func (p *Person) Unmarshal(fromMap map[string]interface{}) error {
	return Unmarshal(fromMap, p)
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

var mapData map[string]interface{}
var structured Person

func init() {
	mapData = map[string]interface{}{
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

func TestMapperBase_Marshal(t *testing.T) {
	result, err := Marshal(structured)
	assert.NoError(t, err)
	assert.Equal(t, mapData, result)
}

func TestMapperBase_Marshal_viaMethod(t *testing.T) {
	result, err := structured.Marshal()
	assert.NoError(t, err)
	assert.Equal(t, mapData, result)
}

func TestMapperBase_Unmarshal_viaMethod(t *testing.T) {
	var result Person
	err := result.Unmarshal(mapData)
	assert.NoError(t, err)
	assert.Equal(t, structured, result)
}

func TestMapperBase_Unmarshal(t *testing.T) {
	var result Person
	err := Unmarshal(mapData, &result)
	assert.NoError(t, err)
	assert.Equal(t, structured, result)
}
