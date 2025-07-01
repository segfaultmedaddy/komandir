package sliceutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFilter(t *testing.T) {
	t.Run("it filters elements", func(t *testing.T) {
		slice := []int{1, 2, 3, 4, 5}
		filtered := Filter(slice, func(i int) bool { return i%2 == 0 })
		expected := []int{2, 4}

		assert.Equal(t, expected, filtered)
	})
}
