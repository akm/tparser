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

	run("String", []rune(`STRING`), &ast.StringType{Name: "STRING"})
	run("ANSI String", []rune(`ANSISTRING`), &ast.StringType{Name: "ANSISTRING"})
	run("Wide String", []rune(`WIDESTRING`), &ast.StringType{Name: "WIDESTRING"})
	run("Short String", []rune(`STRING[100]`), &ast.StringType{Name: "STRING", Length: asttest.NewConstExpr(asttest.NewNumber("100"))})
	run(
		"Short String",
		[]rune(`STRING[ALen + BLen]`),
		&ast.StringType{
			Name: "STRING",
			Length: asttest.NewConstExpr(
				&ast.SimpleExpression{
					Term: asttest.NewTerm(asttest.NewIdent("ALen", asttest.NewIdentLocation(1, 8, 7, 12))),
					AddOpTerms: []*ast.AddOpTerm{
						{AddOp: "+", Term: asttest.NewTerm(asttest.NewIdent("BLen", asttest.NewIdentLocation(1, 15, 14, 19)))},
					},
				},
			),
		},
	)
}
