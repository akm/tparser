package parsertest

import (
	"testing"

	"github.com/akm/tparser/ast"
	"github.com/akm/tparser/ast/asttest"
	"github.com/stretchr/testify/assert"
)

func TestGoal(t *testing.T) {
	text := []rune(`unit U1;
	interface
	implementation
	end.`)

	parser := NewTestParser(&text)
	res, err := parser.ParseGoal()
	if assert.NoError(t, err) {
		asttest.ClearUnitDeclarationMaps(t, res)
		assert.Equal(t, &ast.Unit{
			Ident:                 asttest.NewIdent("U1", asttest.NewIdentLocation(1, 6, 5, 8)),
			InterfaceSection:      &ast.InterfaceSection{},
			ImplementationSection: &ast.ImplementationSection{},
		}, res)
	}
}
