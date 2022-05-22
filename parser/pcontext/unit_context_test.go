package pcontext

import (
	"testing"

	"github.com/akm/tparser/ast"
	"github.com/akm/tparser/ast/astcore"
	"github.com/akm/tparser/ast/asttest"
	"github.com/stretchr/testify/assert"
)

func TestUnitContext(t *testing.T) {
	unit2 := &ast.Unit{
		Ident: asttest.NewIdent("Unit2"),
		InterfaceSection: &ast.InterfaceSection{
			InterfaceDecls: ast.InterfaceDecls{
				ast.TypeSection{
					{Ident: asttest.NewIdent("TType2"), Type: &ast.TypeId{Ident: asttest.NewIdent("String")}},
				},
			},
		},
	}

	usesClauseItemToUnit2 := &ast.UsesClauseItem{
		Ident: asttest.NewIdent("Unit2"),
		Unit:  unit2,
	}

	declMap := astcore.DeclMapImpl{"unit2": usesClauseItemToUnit2.ToDeclarations()[0]}
	assert.NotNil(t, declMap.Get("Unit2"))

	ctx := NewUnitContext(NewProgramContext(), declMap)
	assert.NotNil(t, ctx.Get("Unit2"))
}
