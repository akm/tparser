package parser

import (
	"testing"

	"github.com/akm/tparser/ast"
	"github.com/akm/tparser/ast/asttest"
	"github.com/stretchr/testify/assert"
)

func TestProgram(t *testing.T) {
	run := func(name string, clearLocations bool, text []rune, expected *ast.Program) {
		t.Run(name, func(t *testing.T) {
			parser := NewParser(&text)
			parser.NextToken()
			res, err := parser.ParseProgram()
			if assert.NoError(t, err) {
				if clearLocations {
					asttest.ClearLocations(t, res)
				}
				assert.Equal(t, expected, res)
				// assert.Equal(t, expected.ProgramBlock.Block, res.ProgramBlock.Block)
			}

		})
	}
	run(
		"hello world", false,
		[]rune(`PROGRAM Hello;
begin
  writeln('hello, world');
end.`),
		&ast.Program{
			Ident: asttest.NewIdent("Hello", asttest.NewIdentLocation(1, 9, 8, 14)),
			ProgramBlock: &ast.ProgramBlock{
				Block: &ast.Block{
					CompoundStmt: &ast.CompoundStmt{
						StmtList: ast.StmtList{
							&ast.Statement{
								Body: &ast.CallStatement{
									Designator: asttest.NewDesignator(
										asttest.NewIdent("writeln", asttest.NewIdentLocation(3, 4, 23, 3, 11, 30)),
									),
									ExprList: ast.ExprList{
										asttest.NewExpression(asttest.NewString("'hello, world'")),
									},
								},
							},
						},
					},
				},
			},
		},
	)

	run(
		"hello world with const", true,
		[]rune(`PROGRAM Hello;
const DefaultMessage = 'hello, world';
var msg: string;
begin
  msg := DefaultMessage;
  writeln(msg);
end.
`),
		&ast.Program{
			Ident: asttest.NewIdent("Hello"),
			ProgramBlock: &ast.ProgramBlock{
				Block: &ast.Block{
					DeclSections: ast.DeclSections{
						ast.ConstSection{
							&ast.ConstantDecl{
								Ident:     asttest.NewIdent("DefaultMessage"),
								ConstExpr: asttest.NewConstExpr(asttest.NewString("'hello, world'")),
							},
						},
						ast.VarSection{
							&ast.VarDecl{
								IdentList: asttest.NewIdentList("msg"),
								Type:      asttest.NewStringType("STRING"),
							},
						},
					},
					CompoundStmt: &ast.CompoundStmt{
						StmtList: ast.StmtList{
							&ast.Statement{
								Body: &ast.AssignStatement{
									Designator: asttest.NewDesignator(asttest.NewIdent("msg")),
									Expression: asttest.NewExpression(asttest.NewDesignator("DefaultMessage")),
								},
							},
							&ast.Statement{
								Body: &ast.CallStatement{
									Designator: asttest.NewDesignator(asttest.NewIdent("writeln")),
									ExprList: ast.ExprList{
										asttest.NewExpression(asttest.NewDesignator("msg")),
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
