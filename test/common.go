package test

import (
	"fmt"

	"github.com/madkins23/go-type/data"
)

// PackageName should be set to the known path for this package.
const PackageName = "github.com/madkins23/go-type/test"

//////////////////////////////////////////////////////////////////////////

type Actor interface {
	declaim() string
}

//////////////////////////////////////////////////////////////////////////

type WithExtra struct {
	extra string
}

func (we *WithExtra) Extra() string {
	return we.extra
}

func (we *WithExtra) ClearExtra() {
	we.extra = ""
}

//////////////////////////////////////////////////////////////////////////

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

func (a *Alpha) Marshal() (map[string]interface{}, error) {
	return data.Marshal(a)
}

func (a *Alpha) Unmarshal(fromMap map[string]interface{}) error {
	return data.Unmarshal(fromMap, a)
}

//////////////////////////////////////////////////////////////////////////

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

func (b *Bravo) Marshal() (map[string]interface{}, error) {
	return data.Marshal(b)
}

func (b *Bravo) Unmarshal(fromMap map[string]interface{}) error {
	return data.Unmarshal(fromMap, b)
}
