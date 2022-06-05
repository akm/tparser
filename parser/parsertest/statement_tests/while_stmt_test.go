package statementtest

import (
	"testing"

	"github.com/akm/tparser/ast"
	"github.com/akm/tparser/ast/asttest"
	"github.com/akm/tparser/parser/parsertest"
)

func TestWhileStmt(t *testing.T) {
	parsertest.RunStatementTest(t,
		"one line",
		[]rune(`while Data[I] <> X do I := I + 1;`),
		&ast.Statement{
			Body: &ast.WhileStmt{
				Condition: &ast.Expression{
					SimpleExpression: asttest.NewSimpleExpression(
						&ast.Designator{
							QualId: asttest.NewQualId("Data"),
							Items: ast.DesignatorItems{
								ast.DesignatorItemExprList{
									asttest.NewExpression(asttest.NewQualId("I")),
								},
							},
						},
					),
					RelOpSimpleExpressions: []*ast.RelOpSimpleExpression{
						{RelOp: "<>", SimpleExpression: asttest.NewSimpleExpression(asttest.NewQualId("X"))},
					},
				},
				Statement: &ast.Statement{
					Body: &ast.AssignStatement{
						Designator: asttest.NewDesignator("I"),
						Expression: asttest.NewExpression(
							&ast.SimpleExpression{
								Term: &ast.Term{Factor: asttest.NewDesignatorFactor(asttest.NewIdent("I"))},
								AddOpTerms: []*ast.AddOpTerm{
									{AddOp: "+", Term: asttest.NewTerm(asttest.NewNumber("1"))},
								},
							},
						),
					},
				},
			},
		},
	)

	parsertest.RunStatementTest(t,
		"example1",
		[]rune(`
while I > 0 do
begin
	if Odd(I) then Z := Z * X;
	I := I div 2;
	X := Sqr(X);
end;
`),
		&ast.Statement{
			Body: &ast.WhileStmt{
				Condition: &ast.Expression{
					SimpleExpression: ast.NewSimpleExpression(asttest.NewQualId("I")),
					RelOpSimpleExpressions: []*ast.RelOpSimpleExpression{
						{RelOp: ">", SimpleExpression: asttest.NewSimpleExpression(asttest.NewNumber("0"))},
					},
				},
				Statement: &ast.Statement{
					Body: &ast.CompoundStmt{
						StmtList: ast.StmtList{
							{
								Body: &ast.IfStmt{
									Condition: asttest.NewExpression(
										&ast.DesignatorFactor{
											Designator: asttest.NewDesignator(asttest.NewQualId("Odd")),
											ExprList: ast.ExprList{
												asttest.NewExpression(asttest.NewQualId("I")),
											},
										},
									),
									Then: &ast.Statement{
										Body: &ast.AssignStatement{
											Designator: asttest.NewDesignator("Z"),
											Expression: asttest.NewExpression(
												&ast.Term{
													Factor: asttest.NewDesignatorFactor(asttest.NewIdent("Z")),
													MulOpFactors: []*ast.MulOpFactor{
														{MulOp: "*", Factor: asttest.NewDesignatorFactor(asttest.NewIdent("X"))},
													},
												},
											),
										},
									},
								},
							},
							{
								Body: &ast.AssignStatement{
									Designator: asttest.NewDesignator("I"),
									Expression: asttest.NewExpression(
										&ast.Term{
											Factor: asttest.NewDesignatorFactor(asttest.NewIdent("I")),
											MulOpFactors: []*ast.MulOpFactor{
												{MulOp: "DIV", Factor: asttest.NewNumber("2")},
											},
										},
									),
								},
							},
							{
								Body: &ast.AssignStatement{
									Designator: asttest.NewDesignator("X"),
									Expression: asttest.NewExpression(
										&ast.DesignatorFactor{
											Designator: asttest.NewDesignator(asttest.NewQualId("Sqr")),
											ExprList: ast.ExprList{
												asttest.NewExpression(asttest.NewQualId("X")),
											},
										},
									),
								},
							},
						},
					},
				},
			},
		},
	)

	parsertest.RunStatementTest(t,
		"example2",
		[]rune(`
while not Eof(InputFile) do
begin
	Readln(InputFile, Line);
	Process(Line);
end;
`),
		&ast.Statement{
			Body: &ast.WhileStmt{
				Condition: asttest.NewExpression(
					&ast.Not{
						Factor: &ast.DesignatorFactor{
							Designator: asttest.NewDesignator(asttest.NewQualId("Eof")),
							ExprList: ast.ExprList{
								asttest.NewExpression(asttest.NewQualId("InputFile")),
							},
						},
					},
				),
				Statement: &ast.Statement{
					Body: &ast.CompoundStmt{
						StmtList: ast.StmtList{
							{
								Body: &ast.CallStatement{
									Designator: asttest.NewDesignator(asttest.NewQualId("Readln")),
									ExprList: ast.ExprList{
										asttest.NewExpression(asttest.NewQualId("InputFile")),
										asttest.NewExpression(asttest.NewQualId("Line")),
									},
								},
							},
							{
								Body: &ast.CallStatement{
									Designator: asttest.NewDesignator(asttest.NewQualId("Process")),
									ExprList: ast.ExprList{
										asttest.NewExpression(asttest.NewQualId("Line")),
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
