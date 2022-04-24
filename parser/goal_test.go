package parser

import (
	"testing"

	"github.com/akm/tparser/ast"
	"github.com/stretchr/testify/assert"
)

func TestGoal(t *testing.T) {
	text := []rune(`unit U1;
	interface
	implementation
	end.`)

	parser := NewParser(&text)
	res, err := parser.ParseGoal()
	if assert.NoError(t, err) {
		assert.Equal(t, &ast.Unit{
			Ident:                 *ast.NewIdent("U1"),
			InterfaceSection:      &ast.InterfaceSection{},
			ImplementationSection: &ast.ImplementationSection{},
		}, res)
	}
}
