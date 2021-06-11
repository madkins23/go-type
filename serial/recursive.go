package serial

import (
	"sync"
)

type Recursive interface {
	Open(mapper Mapper)
	Mapper() Mapper
	Close()
}

type WithMapper struct {
	sync.Mutex
	mapper Mapper
}

func (wm *WithMapper) Open(mapper Mapper) {
	wm.Lock()
	wm.mapper = mapper
	wm.Unlock()
}

func (wm *WithMapper) Mapper() Mapper {
	return wm.mapper
}

func (wm *WithMapper) Close() {
	wm.Lock()
	wm.mapper = nil
	wm.Unlock()
}
