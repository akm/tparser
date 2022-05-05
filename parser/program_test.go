package parser

import (
	"testing"

	"github.com/akm/tparser/ast"
	"github.com/akm/tparser/ast/asttest"
	"github.com/stretchr/testify/assert"
)

func TestProgram(t *testing.T) {
	run := func(name string, text []rune, expected *ast.Program) {
		t.Run(name, func(t *testing.T) {
			parser := NewParser(&text)
			parser.NextToken()
			res, err := parser.ParseProgram()
			if assert.NoError(t, err) {
				assert.Equal(t, expected, res)
			}

		})
	}
	run(
		"hello world",
		[]rune(`PROGRAM Hello;
begin
  writeln('hello, world')
end.`),
		&ast.Program{
			Ident: asttest.NewIdent("Hello", asttest.NewIdentLocation(1, 9, 8, 14)),
			ProgramBlock: &ast.ProgramBlock{
				Block: &ast.Block{
					CompoundStmt: &ast.CompoundStmt{
						StmtList: &ast.StmtList{
							Statement: &ast.Statement{
								Body: &ast.CallStatement{
									Designator: asttest.NewDesignator("writeln"),
									ExprList: ast.ExprList{
										asttest.NewExpression(asttest.NewString("hello, world")),
									},
								},
							},
						},
					},
				},
			},
		},
	)
}
