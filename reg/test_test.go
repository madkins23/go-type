package reg

import "fmt"

const packageName = "github.com/madkins23/go-type/reg"

type Stuff interface {
	Info() string
}

type Alpha struct {
	Name   string
	Number float32
}

func (a *Alpha) Info() string {
	return fmt.Sprintf("%s: %f", a.Name, a.Number)
}

type Bravo struct {
	Finished   bool
	Iterations int
}

func (b *Bravo) Info() string {
	return fmt.Sprintf("%t: %d", b.Finished, b.Iterations)
}
