package ast

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnit(t *testing.T) {
	t.Run("Unit implements Goal", func(t *testing.T) {
		var m Goal
		m = &Unit{}
		assert.Implements(t, (*Goal)(nil), m)
	})
	t.Run("InterfaceDecl", func(t *testing.T) {
		assert.Implements(t, (*InterfaceDecl)(nil), &ConstSection{})
		assert.Implements(t, (*InterfaceDecl)(nil), &TypeSection{})
		assert.Implements(t, (*InterfaceDecl)(nil), &VarSection{})
	})
	t.Run("Node", func(t *testing.T) {
		assert.Implements(t, (*Node)(nil), &Unit{})
		assert.Implements(t, (*Node)(nil), &InterfaceSection{})
		assert.Implements(t, (*Node)(nil), &ImplementationSection{})
		assert.Implements(t, (*Node)(nil), &InitSection{})
		assert.Implements(t, (*Node)(nil), &UnitId{})
		assert.Implements(t, (*Node)(nil), &QualId{})
	})
}
