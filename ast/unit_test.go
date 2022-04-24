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
		assert.Implements(t, (*InterfaceDecl)(nil), &ExportedHeading{})
	})
	t.Run("InterfaceSection", func(t *testing.T) {
		assert.Implements(t, (*Node)(nil), &InterfaceSection{})
	})
}
