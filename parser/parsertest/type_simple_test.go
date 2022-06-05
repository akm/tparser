package parsertest

import (
	"testing"

	"github.com/akm/tparser/ast"
	"github.com/akm/tparser/ast/asttest"
)

func TestEnumeratedType(t *testing.T) {
	RunTypeTest(t,
		"card",
		[]rune(`(Club, Diamond, Heart, Spade)`),
		ast.EnumeratedType{
			{Ident: asttest.NewIdent("Club")},
			{Ident: asttest.NewIdent("Diamond")},
			{Ident: asttest.NewIdent("Heart")},
			{Ident: asttest.NewIdent("Spade")},
		},
	)

	RunTypeTest(t,
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
		RunTypeTest(t, name, text, expected, func(tt *BaseTestRunner) {
			tt.ClearLocations = false
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
