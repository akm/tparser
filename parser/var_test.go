package parser

import (
	"testing"

	"github.com/akm/tparser/ast"
	"github.com/akm/tparser/ast/asttest"
	"github.com/stretchr/testify/assert"
)

func TestUnitWithVarSection(t *testing.T) {
	run := func(name string, text []rune, expected *ast.Unit) {
		t.Run(name, func(t *testing.T) {
			parser := NewParser(&text)
			parser.NextToken()
			res, err := parser.ParseUnit()
			if assert.NoError(t, err) {
				asttest.ClearLocations(t, res)
				assert.Equal(t, expected, res)
			}
		})
	}

	run(
		"2 var declarations",
		[]rune(`
		UNIT Unit1;
		INTERFACE
		VAR
		  I: Integer;
		  X, Y: Real;
		IMPLEMENTATION
		END.`),
		&ast.Unit{
			Ident: *asttest.NewIdent("Unit1"),
			InterfaceSection: &ast.InterfaceSection{
				InterfaceDecls: ast.InterfaceDecls{
					ast.VarSection{
						{IdentList: asttest.NewIdentList("I"), Type: &ast.OrdIdent{Name: asttest.NewIdent("Integer")}},
						{IdentList: asttest.NewIdentList("X", "Y"), Type: &ast.RealType{Name: asttest.NewIdent("Real")}},
					},
				},
			},
			ImplementationSection: &ast.ImplementationSection{},
		},
	)
	run(
		"2 var sections",
		[]rune(`
		UNIT Unit1;
		INTERFACE
		VAR
			X, Y, Z: Double;
			I, J, K: Integer;
		VAR
			Digit: 0..9;
			Okay: Boolean;
			A: Integer = 7;
		IMPLEMENTATION
		END.`),
		&ast.Unit{
			Ident: *asttest.NewIdent("Unit1"),
			InterfaceSection: &ast.InterfaceSection{
				InterfaceDecls: []ast.InterfaceDecl{
					ast.VarSection{
						{IdentList: asttest.NewIdentList("X", "Y", "Z"), Type: &ast.RealType{Name: asttest.NewIdent("Double")}},
						{IdentList: asttest.NewIdentList("I", "J", "K"), Type: &ast.OrdIdent{Name: asttest.NewIdent("Integer")}},
					},
					ast.VarSection{
						{IdentList: asttest.NewIdentList("Digit"), Type: &ast.SubrangeType{Low: *asttest.NewConstExpr(asttest.NewNumber("0")), High: *asttest.NewConstExpr(asttest.NewNumber("9"))}},
						{IdentList: asttest.NewIdentList("Okay"), Type: &ast.OrdIdent{Name: asttest.NewIdent("Boolean")}},
						{IdentList: asttest.NewIdentList("A"), Type: &ast.OrdIdent{Name: asttest.NewIdent("Integer")}, ConstExpr: asttest.NewExpression(asttest.NewNumber("7"))},
					},
				},
			},
			ImplementationSection: &ast.ImplementationSection{},
		},
	)
	run(
		"threadvar",
		[]rune(`
		UNIT Unit1;
		INTERFACE
		THREADVAR
			X: Integer;
		IMPLEMENTATION
		END.`),
		&ast.Unit{
			Ident: *asttest.NewIdent("Unit1"),
			InterfaceSection: &ast.InterfaceSection{
				InterfaceDecls: []ast.InterfaceDecl{
					ast.ThreadVarSection{
						{IdentList: asttest.NewIdentList("X"), Type: &ast.OrdIdent{Name: asttest.NewIdent("Integer")}},
					},
				},
			},
			ImplementationSection: &ast.ImplementationSection{},
		},
	)
}

func TestVarSectionl(t *testing.T) {
	run := func(name string, text []rune, expected ast.VarSection) {
		t.Run(name, func(t *testing.T) {
			parser := NewParser(&text)
			parser.NextToken()
			res, err := parser.ParseVarSection()
			if assert.NoError(t, err) {
				assert.Equal(t, expected, res)
			}
		})
	}

	run(
		"Absolute addresses",
		[]rune(`VAR
		Str: string[32];
		StrLen: Byte absolute Str;
		`),
		ast.VarSection{
			&ast.VarDecl{
				IdentList: asttest.NewIdentList(asttest.NewIdent("Str", asttest.NewIdentLocation(2, 4, 6, 7))),
				Type:      asttest.NewStringType("STRING", asttest.NewConstExpr(asttest.NewNumber("32"))),
			},
			&ast.VarDecl{
				IdentList: asttest.NewIdentList(asttest.NewIdent("StrLen", asttest.NewIdentLocation(3, 4, 25, 10))),
				Type:      &ast.OrdIdent{Name: asttest.NewIdent("Byte", asttest.NewIdentLocation(3, 12, 33, 16))},
				Absolute:  asttest.NewVarDeclAbsoluteIdent(asttest.NewIdent("Str", asttest.NewIdentLocation(3, 26, 47, 29))),
			},
		},
	)
	run(
		"With simple ConstExpr",
		[]rune(`VAR A: Integer = 7;`),
		ast.VarSection{
			{
				IdentList: asttest.NewIdentList(asttest.NewIdent("A", asttest.NewIdentLocation(1, 5, 4, 6))),
				Type:      &ast.OrdIdent{Name: asttest.NewIdent("Integer", asttest.NewIdentLocation(1, 8, 7, 15))},
				ConstExpr: asttest.NewExpression(asttest.NewNumber("7")),
			},
		},
	)

	run(
		"var after subrange",
		[]rune(`
		VAR
			Digit: 0..9;
			Okay: Boolean;
		`),
		ast.VarSection{
			{
				IdentList: asttest.NewIdentList(asttest.NewIdent("Digit", asttest.NewIdentLocation(2, 5, 10, 10))),
				Type:      &ast.SubrangeType{Low: *asttest.NewConstExpr(asttest.NewNumber("0")), High: *asttest.NewConstExpr(asttest.NewNumber("9"))}},
			{
				IdentList: asttest.NewIdentList(asttest.NewIdent("Okay", asttest.NewIdentLocation(3, 5, 26, 9))),
				Type:      &ast.OrdIdent{Name: asttest.NewIdent("Boolean", asttest.NewIdentLocation(3, 11, 32, 18))},
			},
		},
	)
}
