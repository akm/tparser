package statementtest

import (
	"testing"

	"github.com/akm/tparser/ast"
	"github.com/akm/tparser/ast/asttest"
	"github.com/akm/tparser/parser/parsertest"
)

func TestTryFinallyStmt(t *testing.T) {
	stateToProgram := func(programName string, statements ...*ast.Statement) *ast.Program {
		return &ast.Program{
			Ident: asttest.NewIdent(programName),
			ProgramBlock: &ast.ProgramBlock{
				Block: &ast.Block{
					Body: &ast.CompoundStmt{
						StmtList: ast.StmtList(statements),
					},
				},
			},
		}
	}

	parsertest.RunProgramTest(t,
		"basic",
		[]rune(`PROGRAM TryFinallyStmtTest;
begin
	Reset(F);
	try
		ProcessFile(F);
	finally
		CloseFile(F);
	end;
end.
`),
		stateToProgram("TryFinallyStmtTest",
			&ast.Statement{
				Body: &ast.CallStatement{
					Designator: asttest.NewDesignator("Reset"),
					ExprList:   ast.ExprList{asttest.NewExpression("F")},
				},
			},
			&ast.Statement{
				Body: &ast.TryFinallyStmt{
					Statements1: ast.StmtList{
						{
							Body: &ast.CallStatement{
								Designator: asttest.NewDesignator("ProcessFile"),
								ExprList:   ast.ExprList{asttest.NewExpression("F")},
							},
						},
					},
					Statements2: ast.StmtList{
						{
							Body: &ast.CallStatement{
								Designator: asttest.NewDesignator("CloseFile"),
								ExprList:   ast.ExprList{asttest.NewExpression("F")},
							},
						},
					},
				},
			},
		),
	)
}
