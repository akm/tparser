package parsertest

import (
	"testing"

	"github.com/akm/tparser/ast"
	"github.com/akm/tparser/ast/asttest"
)

func TestProgram(t *testing.T) {
	run := func(name string, clearLocations bool, text []rune, expected *ast.Program) {
		RunProgramTest(t, name, text, expected, func(tt *BaseTestRunner) {
			tt.ClearLocations = clearLocations
		})
	}

	helloWorldProgram := &ast.Program{
		Ident: asttest.NewIdent("Hello", asttest.NewIdentLocation(1, 9, 8, 14)),
		ProgramBlock: &ast.ProgramBlock{
			Block: &ast.Block{
				Body: &ast.CompoundStmt{
					StmtList: ast.StmtList{
						&ast.Statement{
							Body: &ast.CallStatement{
								Designator: asttest.NewDesignator(
									asttest.NewIdent("writeln", asttest.NewIdentLocation(3, 3, 23, 3, 10, 30)),
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
	}
	run(
		"hello world", false,
		[]rune(`PROGRAM Hello;
begin
  writeln('hello, world');
end.`),
		helloWorldProgram,
	)

	run(
		"hello world omitting semi-colon", false,
		[]rune(`PROGRAM Hello;
begin
  writeln('hello, world')
end.`),
		helloWorldProgram,
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
				Block: func() *ast.Block {
					constDecl := &ast.ConstantDecl{
						Ident:     asttest.NewIdent("DefaultMessage"),
						ConstExpr: asttest.NewConstExpr(asttest.NewString("'hello, world'")),
					}
					varDecl := &ast.VarDecl{
						IdentList: asttest.NewIdentList("msg"),
						Type:      asttest.NewStringType("string"),
					}
					return &ast.Block{
						DeclSections: ast.DeclSections{
							ast.ConstSection{constDecl},
							ast.VarSection{varDecl},
						},
						Body: &ast.CompoundStmt{
							StmtList: ast.StmtList{
								&ast.Statement{
									Body: &ast.AssignStatement{
										Designator: asttest.NewDesignator(
											asttest.NewQualId(asttest.NewIdent("msg"), varDecl.ToDeclarations()[0]),
										),
										Expression: asttest.NewExpression(
											asttest.NewDesignator(
												asttest.NewQualId(asttest.NewIdent("DefaultMessage"), constDecl.ToDeclarations()[0]),
											),
										),
									},
								},
								&ast.Statement{
									Body: &ast.CallStatement{
										Designator: asttest.NewDesignator(asttest.NewIdent("writeln")),
										ExprList: ast.ExprList{
											asttest.NewExpression(
												asttest.NewDesignator(
													asttest.NewQualId(asttest.NewIdent("msg"), varDecl.ToDeclarations()[0]),
												),
											),
										},
									},
								},
							},
						},
					}
				}(),
			},
		},
	)
}
