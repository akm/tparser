package ast

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVar(t *testing.T) {
	t.Run("VarDecl", func(t *testing.T) {
		assert.Implements(t, (*CodeBlock)(nil), &VarDecl{})
	})
	t.Run("ThreadVarDecl", func(t *testing.T) {
		assert.Implements(t, (*CodeBlock)(nil), &ThreadVarDecl{})
	})
	t.Run("VarSection implements InterfaceDecl", func(t *testing.T) {
		var decl InterfaceDecl
		decl = VarSection{}
		assert.Implements(t, (*InterfaceDecl)(nil), decl)
	})
	t.Run("VarDeclAbsoluteConstExpr implements VarDeclAbsolute", func(t *testing.T) {
		var abs VarDeclAbsolute
		abs = &VarDeclAbsoluteConstExpr{}
		assert.Implements(t, (*VarDeclAbsolute)(nil), abs)
	})
	t.Run("ThreadVarSection implements InterfaceDecl", func(t *testing.T) {
		var decl InterfaceDecl
		decl = ThreadVarSection{}
		assert.Implements(t, (*InterfaceDecl)(nil), decl)
	})
}
