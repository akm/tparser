package parser

import (
	"testing"

	"github.com/akm/tparser/ast"
	"github.com/stretchr/testify/assert"
)

func TestUnitWithConstSection(t *testing.T) {
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

	run(
		"const declaration in unit",
		[]rune(`
		UNIT Unit1;
		INTERFACE
		CONST
		  MaxValue = 237;
		  Message1 = 'Out of memory';
		  Max: Integer = 100;
		IMPLEMENTATION
		END.`),
		&ast.Unit{
			Ident: ast.Ident("Unit1"),
			InterfaceSection: &ast.InterfaceSection{
				InterfaceDecls: []ast.InterfaceDecl{
					ast.ConstSection{
						{Ident: ast.Ident("MaxValue"), ConstExpr: *ast.NewExpression(ast.NewNumber("237"))},
						{Ident: ast.Ident("Message1"), ConstExpr: *ast.NewExpression(ast.NewString("'Out of memory'"))},
						{Ident: ast.Ident("Max"), ConstExpr: *ast.NewExpression(ast.NewNumber("100")), Type: &ast.OrdIdent{Name: "Integer"}},
					},
				},
			},
			ImplementationSection: &ast.ImplementationSection{},
		},
	)
}

func TestConstSectionl(t *testing.T) {
	run := func(name string, text []rune, expected ast.ConstSection) {
		t.Run(name, func(t *testing.T) {
			parser := NewParser(&text)
			parser.NextToken()
			res, err := parser.ParseConstSection()
			if assert.NoError(t, err) {
				assert.Equal(t, expected, res)
			}
		})
	}

	run(
		"number const",
		[]rune(`CONST MaxValue = 237;`),
		ast.ConstSection{
			{Ident: ast.Ident("MaxValue"), ConstExpr: *ast.NewExpression(ast.NewNumber("237"))},
		},
	)
	run(
		"number const",
		[]rune(`CONST Max: Integer = 100;`),
		ast.ConstSection{
			{Ident: ast.Ident("Max"), ConstExpr: *ast.NewExpression(ast.NewNumber("100")), Type: &ast.OrdIdent{Name: "Integer"}},
		},
	)
	run(
		"message as identifier",
		[]rune(`CONST Message = 'Out of memory';`),
		ast.ConstSection{
			{Ident: ast.Ident("Message"), ConstExpr: *ast.NewExpression(ast.NewString("'Out of memory'"))},
		},
	)
}
