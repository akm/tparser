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
				asttest.ClearAllRange(res)
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
			Ident: ast.Ident("Unit1"),
			InterfaceSection: &ast.InterfaceSection{
				InterfaceDecls: []ast.InterfaceDecl{
					ast.VarSection{
						{IdentList: ast.IdentList{"I"}, Type: &ast.OrdIdent{Name: ast.Ident("Integer")}},
						{IdentList: ast.IdentList{"X", "Y"}, Type: &ast.RealType{Name: ast.Ident("Real")}},
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
			Ident: ast.Ident("Unit1"),
			InterfaceSection: &ast.InterfaceSection{
				InterfaceDecls: []ast.InterfaceDecl{
					ast.VarSection{
						{IdentList: ast.IdentList{"X", "Y", "Z"}, Type: &ast.RealType{Name: ast.Ident("Double")}},
						{IdentList: ast.IdentList{"I", "J", "K"}, Type: &ast.OrdIdent{Name: ast.Ident("Integer")}},
					},
					ast.VarSection{
						{IdentList: ast.IdentList{"Digit"}, Type: &ast.SubrangeType{Low: *ast.NewConstExpr(ast.NewNumber("0")), High: *ast.NewConstExpr(ast.NewNumber("9"))}},
						{IdentList: ast.IdentList{"Okay"}, Type: &ast.OrdIdent{Name: ast.Ident("Boolean")}},
						{IdentList: ast.IdentList{"A"}, Type: &ast.OrdIdent{Name: ast.Ident("Integer")}, ConstExpr: ast.NewExpression(ast.NewNumber("7"))},
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
			Ident: ast.Ident("Unit1"),
			InterfaceSection: &ast.InterfaceSection{
				InterfaceDecls: []ast.InterfaceDecl{
					ast.ThreadVarSection{
						{IdentList: ast.IdentList{"X"}, Type: &ast.OrdIdent{Name: ast.Ident("Integer")}},
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
				IdentList: ast.IdentList{"Str"},
				Type:      &ast.StringType{Name: "STRING", Length: ast.NewConstExpr(ast.NewNumber("32"))},
				CodeBlockNode: *asttest.NewCodeBlockNode(
					asttest.CodePosition(6, 2, 4),
					asttest.CodePosition(21, 2, 19),
				),
			},
			&ast.VarDecl{
				IdentList: ast.IdentList{"StrLen"},
				Type:      &ast.OrdIdent{Name: ast.Ident("Byte")},
				Absolute:  ast.NewVarDeclAbsoluteIdent("Str"),
				CodeBlockNode: *asttest.NewCodeBlockNode(
					asttest.CodePosition(25, 3, 4),
					asttest.CodePosition(50, 3, 29),
				),
			},
		},
	)
	run(
		"With simple ConstExpr",
		[]rune(`VAR A: Integer = 7;`),
		ast.VarSection{
			{
				IdentList: ast.IdentList{"A"},
				Type:      &ast.OrdIdent{Name: ast.Ident("Integer")}, ConstExpr: ast.NewExpression(ast.NewNumber("7")),
				CodeBlockNode: *asttest.NewCodeBlockNode(
					asttest.CodePosition(4, 1, 5),
					asttest.CodePosition(18, 1, 19),
				),
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
				IdentList: ast.IdentList{"Digit"},
				Type:      &ast.SubrangeType{Low: *ast.NewConstExpr(ast.NewNumber("0")), High: *ast.NewConstExpr(ast.NewNumber("9"))},
				CodeBlockNode: *asttest.NewCodeBlockNode(
					asttest.CodePosition(10, 2, 5),
					asttest.CodePosition(21, 2, 16),
				),
			},
			{
				IdentList: ast.IdentList{"Okay"},
				Type:      &ast.OrdIdent{Name: ast.Ident("Boolean")},
				CodeBlockNode: *asttest.NewCodeBlockNode(
					asttest.CodePosition(26, 3, 5),
					asttest.CodePosition(39, 3, 18),
				),
			},
		},
	)

}

func TestThreadVarSectionl(t *testing.T) {
	run := func(name string, text []rune, expected ast.ThreadVarSection) {
		t.Run(name, func(t *testing.T) {
			parser := NewParser(&text)
			parser.NextToken()
			res, err := parser.ParseThreadVarSection()
			if assert.NoError(t, err) {
				assert.Equal(t, expected, res)
			}
		})
	}
	run(
		"With simple ConstExpr",
		[]rune(`THREADVAR X: Integer;`),
		ast.ThreadVarSection{
			{
				IdentList: ast.IdentList{"X"},
				Type:      &ast.OrdIdent{Name: ast.Ident("Integer")},
				CodeBlockNode: *asttest.NewCodeBlockNode(
					asttest.CodePosition(10, 1, 11),
					asttest.CodePosition(20, 1, 21),
				),
			},
		},
	)
}
