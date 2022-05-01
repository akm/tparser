package ast

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUsesClause(t *testing.T) {
	t.Run("Node", func(t *testing.T) {
		assert.Implements(t, (*Node)(nil), UsesClause{})
		assert.Implements(t, (*Node)(nil), &UnitRef{})
	})
}
