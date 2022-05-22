package parsertest

import (
	"testing"

	"github.com/akm/tparser/ast"
	"github.com/akm/tparser/ast/asttest"
	"github.com/stretchr/testify/assert"
)

func TestEnumeratedType(t *testing.T) {
	run := func(name string, text []rune, expected ast.Type) {
		t.Run(name, func(t *testing.T) {
			parser := NewTestParser(&text)
			parser.NextToken()
			res, err := parser.ParseType()
			if assert.NoError(t, err) {
				asttest.ClearLocations(t, res)
				assert.Equal(t, expected, res)
			}
		})
	}

	run(
		"card",
		[]rune(`(Club, Diamond, Heart, Spade)`),
		ast.EnumeratedType{
			{Ident: asttest.NewIdent("Club")},
			{Ident: asttest.NewIdent("Diamond")},
			{Ident: asttest.NewIdent("Heart")},
			{Ident: asttest.NewIdent("Spade")},
		},
	)

	run(
		"Enumerated types with explicitly assigned ordinality",
		[]rune(`(Small = 5, Medium = 10, Large = Small + Medium)`),
		func() ast.EnumeratedType {
			small := &ast.EnumeratedTypeElement{Ident: asttest.NewIdent("Small"), ConstExpr: asttest.NewConstExpr(asttest.NewNumber("5"))}
			medium := &ast.EnumeratedTypeElement{Ident: asttest.NewIdent("Medium"), ConstExpr: asttest.NewConstExpr(asttest.NewNumber("10"))}
			return ast.EnumeratedType{
				small,
				medium,
				{Ident: asttest.NewIdent("Large"), ConstExpr: asttest.NewConstExpr(
					&ast.SimpleExpression{
						Term: asttest.NewTerm(asttest.NewQualId("Small", small.ToDeclarations()[0])),
						AddOpTerms: []*ast.AddOpTerm{
							{AddOp: "+", Term: asttest.NewTerm(asttest.NewQualId("Medium", medium.ToDeclarations()[0]))},
						},
					},
				)},
			}
		}(),
	)
}

func TestSubrangeType(t *testing.T) {
	run := func(name string, text []rune, expected ast.Type) {
		t.Run(name, func(t *testing.T) {
			parser := NewTestParser(&text)
			parser.NextToken()
			res, err := parser.ParseType()
			if assert.NoError(t, err) {
				if !assert.Equal(t, expected, res) {
					asttest.AssertType(t, expected, res)
				}
			}
		})
	}

	run(
		"subrange type of enumerated type",
		[]rune(`Green..White`),
		&ast.SubrangeType{
			Low:  asttest.NewConstExpr(asttest.NewIdent("Green", asttest.NewIdentLocation(1, 1, 0, 6))),
			High: asttest.NewConstExpr(asttest.NewIdent("White", asttest.NewIdentLocation(1, 8, 7, 1, 13, 12))),
		},
	)
	run(
		"subrange type of number",
		[]rune(`-128..127`),
		&ast.SubrangeType{Low: asttest.NewConstExpr(asttest.NewNumber("-128")), High: asttest.NewConstExpr(asttest.NewNumber("127"))},
	)
	run(
		"subrange type of character",
		[]rune(`'A'..'Z'`),
		&ast.SubrangeType{Low: asttest.NewConstExpr(asttest.NewString("'A'")), High: asttest.NewConstExpr(asttest.NewString("'Z'"))},
	)
}
