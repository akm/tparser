package parser

import (
	"testing"

	"github.com/akm/tparser/ast"
	"github.com/stretchr/testify/assert"
)

func TestUnit(t *testing.T) {
	type pattern struct {
		name     string
		text     []rune
		expected *ast.Unit
	}

	patterns := []*pattern{
		{
			// "unit" is loaded already in goal.go
			name: "simplest unit",
			text: []rune(`U1;
			interface
			implementation
			end.`),
			expected: &ast.Unit{
				Ident:                 ast.Ident("U1"),
				InterfaceSection:      &ast.InterfaceSection{},
				ImplementationSection: &ast.ImplementationSection{},
			},
		},
	}

	for _, ptn := range patterns {
		parser := NewParser(&ptn.text)
		res, err := parser.ParseUnit()
		if assert.NoError(t, err) {
			assert.Equal(t, ptn.expected, res)
		}
	}
}

func TestInterfaceSection(t *testing.T) {
	type pattern struct {
		name     string
		text     []rune
		expected *ast.InterfaceSection
	}

	patterns := []*pattern{
		{
			// "unit" is loaded already in goal.go
			name: "Uses only",
			text: []rune(`USES U1,U2,U3;`),
			expected: &ast.InterfaceSection{
				UsesClause: &ast.UsesClause{"U1", "U2", "U3"},
			},
		},
	}

	for _, ptn := range patterns {
		parser := NewParser(&ptn.text)
		res, err := parser.ParseInterfaceSection()
		if assert.NoError(t, err) {
			assert.Equal(t, ptn.expected, res)
		}
	}
}
