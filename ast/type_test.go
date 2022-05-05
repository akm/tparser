package ast

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestType(t *testing.T) {
	t.Run("Node", func(t *testing.T) {
		assert.Implements(t, (*Node)(nil), TypeSection{})
		assert.Implements(t, (*Node)(nil), &TypeDecl{})
		assert.Implements(t, (*Node)(nil), &TypeId{})
	})
}
