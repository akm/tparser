package ast

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGoal(t *testing.T) {
	t.Run("Goal", func(t *testing.T) {
		// assert.Implements(t, (*Goal)(nil), &Program{})
		// assert.Implements(t, (*Goal)(nil), &Package{})
		// assert.Implements(t, (*Goal)(nil), &Library{})
		assert.Implements(t, (*Goal)(nil), &Unit{})
	})
}
