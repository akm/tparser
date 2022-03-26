package parser

import (
	"testing"

	"github.com/akm/tparser/ast"
	"github.com/stretchr/testify/assert"
)

func TestEnumeratedType(t *testing.T) {
	run := func(name string, text []rune, expected ast.Type) {
		t.Run(name, func(t *testing.T) {
			parser := NewParser(&text)
			parser.NextToken()
			res, err := parser.ParseType()
			if assert.NoError(t, err) {
				assert.Equal(t, expected, res)
			}
		})
	}

	run(
		"card",
		[]rune(`(Club, Diamond, Heart, Spade)`),
		ast.EnumeratedType{
			{Ident: ast.Ident("Club")},
			{Ident: ast.Ident("Diamond")},
			{Ident: ast.Ident("Heart")},
			{Ident: ast.Ident("Spade")},
		},
	)

	run(
		"Enumerated types with explicitly assigned ordinality",
		[]rune(`(Small = 5, Medium = 10, Large = Small + Medium)`),
		ast.EnumeratedType{
			{Ident: ast.Ident("Small"), ConstExpr: ast.NewConstExpr(ast.NewNumber("5"))},
			{Ident: ast.Ident("Medium"), ConstExpr: ast.NewConstExpr(ast.NewNumber("10"))},
			{Ident: ast.Ident("Large"), ConstExpr: ast.NewConstExpr(
				&ast.SimpleExpression{
					Term: *ast.NewTerm("Small"),
					AddOpTerms: []*ast.AddOpTerm{
						{AddOp: "+", Term: *ast.NewTerm("Medium")},
					},
				},
			)},
		},
	)

}

func TestSubrangeType(t *testing.T) {
	run := func(name string, text []rune, expected ast.Type) {
		t.Run(name, func(t *testing.T) {
			parser := NewParser(&text)
			parser.NextToken()
			res, err := parser.ParseType()
			if assert.NoError(t, err) {
				assert.Equal(t, expected, res)
			}
		})
	}

	run(
		"subrange type of enumerated type",
		[]rune(`Green..White`),
		&ast.SubrangeType{Low: *ast.NewConstExpr("Green"), High: *ast.NewConstExpr("White")},
	)
	run(
		"subrange type of number",
		[]rune(`-128..127`),
		&ast.SubrangeType{Low: *ast.NewConstExpr(ast.NewNumber("-128")), High: *ast.NewConstExpr(ast.NewNumber("127"))},
	)
	run(
		"subrange type of character",
		[]rune(`'A'..'Z'`),
		&ast.SubrangeType{Low: *ast.NewConstExpr(ast.NewString("'A'")), High: *ast.NewConstExpr(ast.NewString("'Z'"))},
	)
}
