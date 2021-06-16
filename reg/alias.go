package reg

import (
	"fmt"
	"sync"
)

// AddAlias provides a package-specific alias and Registry combination.
// This simplifies registration of types in a package with a common alias.
type Alias struct {
	Registry
	alias    string
	aliased  bool
	updating sync.Mutex
}

// NewAlias returns a package-specific Registry with the given alias.
// If the provided registry is nil a non-current Registry will be provided.
func NewAlias(alias string, registry Registry) *Alias {
	if registry == nil {
		registry = theOne
	}
	return &Alias{
		alias:    alias,
		Registry: registry,
	}
}

// Register the type for the specified example object.
// Generates the embedded Registry.Alias() call with first use.
// Actual registration passed along to package registry object.
func (a *Alias) Register(example interface{}) error {
	if !a.aliased {
		a.updating.Lock()
		if !a.aliased {
			if err := a.AddAlias(a.alias, example); err != nil {
				return fmt.Errorf("register alias: %w", err)
			}
			a.aliased = true
		}
		a.updating.Unlock()
	}

	if err := a.Registry.Register(example); err != nil {
		return fmt.Errorf("register example: %w", err)
	}

	return nil
}
