package ast

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConst(t *testing.T) {
	t.Run("Node", func(t *testing.T) {
		assert.Implements(t, (*Node)(nil), ConstSection{})
		assert.Implements(t, (*DeclSection)(nil), ConstSection{})
		assert.Implements(t, (*InterfaceDecl)(nil), ConstSection{})
		assert.Implements(t, (*Node)(nil), &ConstantDecl{})
		assert.Implements(t, (*Node)(nil), &ConstExpr{})
	})
}
