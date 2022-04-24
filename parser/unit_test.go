package parser

import (
	"testing"

	"github.com/akm/tparser/ast"
	"github.com/stretchr/testify/assert"
)

func TestUnit(t *testing.T) {
	run := func(name string, text []rune, expected *ast.Unit) {
		t.Run(name, func(t *testing.T) {
			parser := NewParser(&text)
			parser.NextToken()
			res, err := parser.ParseUnit()
			if assert.NoError(t, err) {
				assert.Equal(t, expected, res)
			}

		})
	}
	// "unit" is loaded already in goal.go
	run(
		"simplest unit",
		[]rune(`UNIT U1;
			interface
			implementation
			end.`),
		&ast.Unit{
			Ident:                 *ast.NewIdent("U1"),
			InterfaceSection:      &ast.InterfaceSection{},
			ImplementationSection: &ast.ImplementationSection{},
		},
	)
}

func TestInterfaceSection(t *testing.T) {
	run := func(name string, text []rune, expected *ast.InterfaceSection) {
		t.Run(name, func(t *testing.T) {
			parser := NewParser(&text)
			parser.NextToken()
			res, err := parser.ParseInterfaceSection()
			if assert.NoError(t, err) {
				assert.Equal(t, expected, res)
			}

		})
	}

	// "unit" is loaded already in goal.go
	run(
		"Uses only",
		[]rune(`INTERFACE USES U1,U2,U3;`),
		&ast.InterfaceSection{
			UsesClause: &ast.UsesClause{
				ast.NewUnitRef("U1"),
				ast.NewUnitRef("U2"),
				ast.NewUnitRef("U3"),
			},
		},
	)
}
