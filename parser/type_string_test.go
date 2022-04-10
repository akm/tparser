package parser

import (
	"testing"

	"github.com/akm/tparser/ast"
	"github.com/stretchr/testify/assert"
)

func TestStringType(t *testing.T) {
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

	run("String", []rune(`STRING`), &ast.StringType{Name: "STRING"})
	run("ANSI String", []rune(`ANSISTRING`), &ast.StringType{Name: "ANSISTRING"})
	run("Wide String", []rune(`WIDESTRING`), &ast.StringType{Name: "WIDESTRING"})
	run("Short String", []rune(`STRING[100]`), &ast.StringType{Name: "STRING", Length: ast.NewConstExpr(ast.NewNumber("100"))})
	run(
		"Short String",
		[]rune(`STRING[ALen + BLen]`),
		&ast.StringType{
			Name: "STRING",
			Length: ast.NewConstExpr(
				&ast.SimpleExpression{
					Term: *ast.NewTerm("ALen"),
					AddOpTerms: []*ast.AddOpTerm{
						{AddOp: "+", Term: *ast.NewTerm("BLen")},
					},
				},
			),
		},
	)
}