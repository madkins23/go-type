package test

import (
	"encoding/json"
	"fmt"
)

// PackageName should be set to the known path for this package.
const PackageName = "github.com/madkins23/go-type/test"

//////////////////////////////////////////////////////////////////////////

type Actor interface {
	declaim() string
}

type WithExtra struct {
	extra string
}

func (we *WithExtra) Extra() string {
	return we.extra
}

func (we *WithExtra) ClearExtra() {
	we.extra = ""
}

type Alpha struct {
	Name    string
	Percent float32 `yaml:"percentDone"`
	WithExtra
}

func NewAlpha(name string, percent float32, extra string) *Alpha {
	a := &Alpha{Name: name, Percent: percent}
	a.extra = extra
	return a
}

func (a *Alpha) declaim() string {
	return fmt.Sprintf("%s is %6.2f%%  complete", a.Name, a.Percent)
}

type Bravo struct {
	Finished   bool
	Iterations int
	WithExtra
}

func NewBravo(finished bool, iterations int, extra string) *Bravo {
	b := &Bravo{Finished: finished, Iterations: iterations}
	b.extra = extra
	return b
}

func (b *Bravo) declaim() string {
	var finished string
	if !b.Finished {
		finished = "not "
	}
	return fmt.Sprintf("%sfinished after %d iterations", finished, b.Iterations)
}

//////////////////////////////////////////////////////////////////////////
// Configure RegistryItem methods in a way that allows behavior to be defined in test scripts.

var CopyMapFromItemFn func(toMap map[string]interface{}, fromItem interface{}) error
var CopyItemFromMapFn func(toItem interface{}, fromMap map[string]interface{}) error

var errCopyFnMissing = fmt.Errorf("no copy function")

func (a *Alpha) PushToMap(toMap map[string]interface{}) error {
	if CopyMapFromItemFn == nil {
		return errCopyFnMissing
	}

	return CopyMapFromItemFn(toMap, a)
}

func (a *Alpha) PullFromMap(fromMap map[string]interface{}) error {
	if CopyItemFromMapFn == nil {
		return errCopyFnMissing
	}

	return CopyItemFromMapFn(a, fromMap)
}

func (b *Bravo) PushToMap(toMap map[string]interface{}) error {
	if CopyMapFromItemFn == nil {
		return errCopyFnMissing
	}

	return CopyMapFromItemFn(toMap, b)
}

func (b *Bravo) PullFromMap(fromMap map[string]interface{}) error {
	if CopyItemFromMapFn == nil {
		return errCopyFnMissing
	}

	return CopyItemFromMapFn(b, fromMap)
}

// CopyMapFromItemJSON is a default copy mechanism for testing.
func CopyMapFromItemJSON(toMap map[string]interface{}, fromItem interface{}) error {
	if bytes, err := json.Marshal(fromItem); err != nil {
		return fmt.Errorf("marshaling from %v: %w", fromItem, err)
	} else if err = json.Unmarshal(bytes, &toMap); err != nil {
		return fmt.Errorf("unmarshaling to %v: %w", toMap, err)
	}

	return nil
}

// CopyItemFromMapJSON is a default copy mechanism for testing.
func CopyItemFromMapJSON(toItem interface{}, fromMap map[string]interface{}) error {
	if bytes, err := json.Marshal(fromMap); err != nil {
		return fmt.Errorf("marshaling from %v: %w", fromMap, err)
	} else if err = json.Unmarshal(bytes, toItem); err != nil {
		return fmt.Errorf("unmarshaling to %v: %w", toItem, err)
	}

	return nil
}
