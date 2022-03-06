package parser

import (
	"testing"

	"github.com/akm/opparser/ast"
	"github.com/stretchr/testify/assert"
)

func TestUnit(t *testing.T) {
	// "unit" is loaded already in goal.go
	text := []rune(`U1;
	interface
	implementation
	end.`)

	parser := NewParser(&text)
	res, err := parser.ParseUnit()
	if assert.NoError(t, err) {
		assert.Equal(t, &ast.Unit{
			Ident:                 ast.Ident("U1"),
			InterfaceSection:      &ast.InterfaceSection{},
			ImplementationSection: &ast.ImplementationSection{},
		}, res)
	}
}