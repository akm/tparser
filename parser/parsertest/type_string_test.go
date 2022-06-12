package parsertest

import (
	"testing"

	"github.com/akm/tparser/ast"
	"github.com/akm/tparser/ast/asttest"
)

func TestStringType(t *testing.T) {
	run := func(name string, text []rune, expected ast.Type) {
		RunTypeTest(t, name, text, expected, func(tt *BaseTestRunner) {
			tt.ClearLocations = false
		})
	}

	run("String", []rune(`STRING`), asttest.NewStringType(asttest.NewIdent("STRING", asttest.NewIdentLocation(1, 1, 0, 7))))
	run("ANSI String", []rune(`ANSISTRING`), asttest.NewStringType(asttest.NewIdent("ANSISTRING", asttest.NewIdentLocation(1, 1, 0, 11))))
	run("Wide String", []rune(`WIDESTRING`), asttest.NewStringType(asttest.NewIdent("WIDESTRING", asttest.NewIdentLocation(1, 1, 0, 11))))
	run("Short String", []rune(`STRING[100]`), asttest.NewFixedStringType(
		asttest.NewIdent("STRING", asttest.NewIdentLocation(1, 1, 0, 7)),
		asttest.NewConstExpr(asttest.NewNumber("100"))),
	)
	run(
		"Short String",
		[]rune(`STRING[ALen + BLen]`),
		asttest.NewFixedStringType(
			asttest.NewIdent("STRING", asttest.NewIdentLocation(1, 1, 0, 7)),
			asttest.NewConstExpr(
				&ast.SimpleExpression{
					Term: asttest.NewTerm(asttest.NewIdent("ALen", asttest.NewIdentLocation(1, 8, 7, 12))),
					AddOpTerms: []*ast.AddOpTerm{
						{AddOp: "+", Term: asttest.NewTerm(asttest.NewIdent("BLen", asttest.NewIdentLocation(1, 15, 14, 19)))},
					},
				},
			),
		),
	)
}
