package reg

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"unicode"
	"unicode/utf8"
)

// Registry is the type registry interface.
// A type registry tracks specific types by name, a facility not native to Go.
// A type name in the registry is made up of package path and local type name.
// Aliases may be specified to shorten the path to manageable lengths.
//
// The methods on this object are duplicated as top-level functions which
// invoke the predefined global object.
type Registry interface {
	// AddAlias creates an alias to be used to shorten names.
	// The alias must exist prior to registering applicable types.
	// Redefining a pre-existing alias is an error.
	AddAlias(alias string, example interface{}) error

	// Register a type by providing an example object.
	Register(example interface{}) error

	// Make creates a new instance of the example object with the specified name.
	// The new instance will be created with fields filled with zero values.
	Make(name string) (interface{}, error)

	// NameFor returns the current name for the registered type of the specified object.
	NameFor(item interface{}) (string, error)

	// Clear removes all previous aliases and registrations.
	// Intended for use in unit tests in the same package to avoid overlaps.
	Clear()
}

// NewRegistry creates a new Registry object of the default internal type.
func NewRegistry() Registry {
	return &registry{
		aliases: make(map[string]string),
		byName:  make(map[string]*registration),
		byType:  make(map[reflect.Type]*registration),
	}
}

//////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////

// Make sure the interface is satisfied at compile time.
var _ Registry = &registry{}

// Default Registry implementation.
type registry struct {
	// byName supports lookup of registrations by 'name'.
	// Full names and aliases are both entered herein.
	byName map[string]*registration

	// byType supports lookup of registrations by type.
	byType map[reflect.Type]*registration

	// alias maps shortened 'alias' strings to path prefix to shorten names.
	aliases map[string]string
}

// Registration structure groups data from indexes.
type registration struct {
	// currentName includes package path and type name.
	currentName string

	// allNames is the set of all possible type names (i.e. including aliased).
	// The best one will always be in currentName.
	allNames []string

	// typeObj is the reflect.Type object for the example object.
	typeObj reflect.Type
}

//////////////////////////////////////////////////////////////////////////

// AddAlias creates an alias to be used to shorten names.
// The alias must exist prior to registering applicable types.
// Redefining a pre-existing alias is an error.
func (reg *registry) AddAlias(alias string, example interface{}) error {
	if _, found := reg.aliases[alias]; found {
		return fmt.Errorf("can't redefine alias %s", alias)
	}

	exampleType := reflect.TypeOf(example)
	if exampleType == nil {
		return fmt.Errorf("find type for alias %s (%v)", alias, example)
	}

	if exampleType.Kind() == reflect.Ptr {
		exampleType = exampleType.Elem()
		if exampleType == nil {
			return fmt.Errorf("no elem type for alias %s (%v)", alias, example)
		}
	}

	pkgPath := exampleType.PkgPath()
	if pkgPath == "" {
		return fmt.Errorf("no package path for alias %s (%v)", alias, example)
	}

	reg.aliases[alias] = pkgPath
	return nil
}

// Register a type by providing an example object.
func (reg *registry) Register(example interface{}) error {
	// Get reflected type for example object.
	exType := reflect.TypeOf(example)
	if exType != nil && exType.Kind() == reflect.Ptr {
		exType = exType.Elem()
	}
	if exType == nil {
		return fmt.Errorf("no reflected type for %v", example)
	}

	// Check for previous record.
	if _, ok := reg.byType[exType]; ok {
		return fmt.Errorf("previous registration for type %v", exType)
	}

	// Get type name without any pointer asterisks.
	typeName := exType.String()
	if strings.HasPrefix(typeName, "*") {
		typeName = strings.TrimLeft(typeName, "*")
	}

	typeNameSplit := strings.Split(typeName, ".")
	r, _ := utf8.DecodeRuneInString(typeNameSplit[len(typeNameSplit)-1])
	if !unicode.IsUpper(r) {
		return fmt.Errorf("type '%s' is private", typeName)
	}

	// Create registration record for this type.
	item := &registration{
		currentName: typeName,
		allNames:    make([]string, 1, len(reg.aliases)+1),
		typeObj:     exType,
	}

	// Initialize default name to full name with package and type.
	name, aliases, err := reg.genNames(example, true)
	if err != nil {
		return fmt.Errorf("getting type name of example: %w", err)
	}

	item.currentName = name
	item.allNames[0] = name
	for _, alias := range aliases {
		item.allNames = append(item.allNames, alias)
	}

	// Add name lookups for all default and aliased names.
	reg.byName[name] = item
	for _, name := range item.allNames {
		reg.byName[name] = item
	}

	// Add type lookup.
	reg.byType[exType] = item

	return nil
}

var errItemIsNil = errors.New("item is nil")

// NameFor returns the current name for the registered type of the specified object.
func (reg *registry) NameFor(item interface{}) (string, error) {
	itemType := reflect.TypeOf(item)
	if itemType == nil {
		return "", errItemIsNil
	}
	if itemType.Kind() == reflect.Ptr {
		itemType = itemType.Elem()
	}

	registration, ok := reg.byType[itemType]
	if !ok {
		return "", fmt.Errorf("no registration for type %s", itemType)
	}

	return registration.currentName, nil
}

// Make creates a new instance of the example object with the specified name.
// The new instance will be created with fields filled with zero values.
func (reg *registry) Make(name string) (interface{}, error) {
	item, found := reg.byName[name]
	if !found {
		return nil, fmt.Errorf("no registration for type named '%s'", name)
	}

	return reflect.New(item.typeObj).Interface(), nil
}

func (reg *registry) Clear() {
	reg.aliases = make(map[string]string)
	reg.byName = make(map[string]*registration)
	reg.byType = make(map[reflect.Type]*registration)
}

//////////////////////////////////////////////////////////////////////////

func genNameFromInterface(example interface{}) (string, error) {
	itemType := reflect.TypeOf(example)
	if itemType == nil {
		return "", fmt.Errorf("no type for item %v", example)
	}

	if itemType.Kind() == reflect.Ptr {
		itemType = itemType.Elem()
	}

	path := itemType.PkgPath()
	if path == "" {
		return "", fmt.Errorf("no path for type %s", itemType.Name())
	}

	last := strings.LastIndex(path, "/")
	if last < 0 {
		return "", fmt.Errorf("no slash in %s", path)
	}

	final := path[last:]
	name := itemType.Name()

	if strings.HasPrefix(name, final+".") {
		name = name[len(final)+1:]
	}

	return path + "/" + name, nil
}

// genNames creates the possible names for the type represented by the example object.
// Returns the 'canonical' name, an optional array of aliased names per current aliases, and any error.
// If the aliased argument is true a possibly empty array will be returned for the second argument otherwise nil.
func (reg *registry) genNames(example interface{}, aliased bool) (string, []string, error) {
	// Initialize default name to full name with package and type.
	name, err := genNameFromInterface(example)
	if err != nil {
		return "", nil, fmt.Errorf("generating basic name: %w", err)
	}

	var aliases []string
	if aliased {
		aliases = make([]string, 0, len(reg.aliases))

		// Look for any possible aliases for the type and add them to the list of all names.
		for alias, prefixPath := range reg.aliases {
			if strings.HasPrefix(name, prefixPath) {
				aliases = append(aliases, "["+alias+"]"+name[len(prefixPath)+1:])
			}
		}

		// Choose default name again from shortest, therefore most likely an aliased name if there are any.
		nameLen := len(name)
		for _, alias := range aliases {
			// Using <= favors later aliases of same size.
			if len(alias) <= nameLen {
				name = alias
			}
		}
	}

	return name, aliases, nil
}
