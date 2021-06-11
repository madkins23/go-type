package serial

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWithMapperIsRecursive(t *testing.T) {
	var wm interface{} = &WithMapper{}
	rec, ok := wm.(Recursive)
	assert.True(t, ok)
	assert.NotNil(t, rec)
}
