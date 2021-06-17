package data

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"
)

type Marshalable interface {
	Marshal() (Map, error)
	Unmarshal(fromMap Map) error
}

type Person struct {
	Name   string
	Age    int
	Emails []string
	Extra  map[string]string
}

func (p *Person) Marshal() (Map, error) {
	return Marshal(p)
}

func (p *Person) Unmarshal(fromMap Map) error {
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

var mapData Map
var structured Person
var iface Marshalable

func init() {
	mapData = Map{
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
	iface = &Person{
		Name:   name,
		Age:    age,
		Emails: emails(),
		Extra:  extra(),
	}
}

func TestMapperBase_Marshal(t *testing.T) {
	result, err := Marshal(structured)
	require.NoError(t, err)
	assert.Equal(t, mapData, result)
	result, err = Marshal(iface)
	require.NoError(t, err)
	assert.Equal(t, mapData, result)
}

func TestMapperBase_Marshal_viaMethod(t *testing.T) {
	result, err := structured.Marshal()
	require.NoError(t, err)
	assert.Equal(t, mapData, result)
	result, err = iface.Marshal()
	require.NoError(t, err)
	assert.Equal(t, mapData, result)
}

func TestMapperBase_Unmarshal_viaMethod(t *testing.T) {
	var result1 Person
	err := result1.Unmarshal(mapData)
	require.NoError(t, err)
	assert.Equal(t, structured, result1)
	var result2 Person
	require.NoError(t, result2.Unmarshal(mapData))
}

func TestMapperBase_Unmarshal(t *testing.T) {
	var result1 Person
	err := Unmarshal(mapData, &result1)
	require.NoError(t, err)
	assert.Equal(t, structured, result1)
	var result2 Person
	var ires interface{} = &result2
	err = Unmarshal(mapData, ires)
	require.NoError(t, err)
	assert.Equal(t, structured, result2)
}
