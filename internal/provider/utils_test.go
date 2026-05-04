package provider

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_sliceToString(t *testing.T) {
	foo := []any{"test"}
	bar := sliceToString(foo)
	assert.Equal(t, bar, []string{"test"})
}
