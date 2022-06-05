package statementtest

import (
	"testing"

	"github.com/akm/tparser/ast"
	"github.com/akm/tparser/ast/asttest"
	"github.com/akm/tparser/parser/parsertest"
)

func TestRepeatStmt(t *testing.T) {
	parsertest.RunStatementTest(t,
		"simple1",
		[]rune(`
repeat
	K := I mod J;
	I := J;
	J := K;
until J = 0;
`),
		&ast.Statement{
			Body: &ast.RepeatStmt{
				StmtList: ast.StmtList{
					{
						Body: &ast.AssignStatement{
							Designator: asttest.NewDesignator("K"),
							Expression: asttest.NewExpression(
								&ast.Term{
									Factor: asttest.NewDesignatorFactor(asttest.NewQualId("I")),
									MulOpFactors: []*ast.MulOpFactor{
										{MulOp: "MOD", Factor: asttest.NewDesignatorFactor(asttest.NewQualId("J"))},
									},
								},
							),
						},
					},
					{
						Body: &ast.AssignStatement{
							Designator: asttest.NewDesignator("I"),
							Expression: asttest.NewExpression(asttest.NewQualId("J")),
						},
					},
					{
						Body: &ast.AssignStatement{
							Designator: asttest.NewDesignator("J"),
							Expression: asttest.NewExpression(asttest.NewQualId("K")),
						},
					},
				},
				Condition: &ast.Expression{
					SimpleExpression: asttest.NewSimpleExpression(asttest.NewQualId("J")),
					RelOpSimpleExpressions: []*ast.RelOpSimpleExpression{
						{RelOp: "=", SimpleExpression: asttest.NewSimpleExpression(asttest.NewNumber("0"))},
					},
				},
			},
		},
	)

	parsertest.RunStatementTest(t,
		"simple1",
		[]rune(`
repeat
	Write('Enter a value (0..9): ');
	Readln(I);
until (I >= 0) and (I <= 9);
`),
		&ast.Statement{
			Body: &ast.RepeatStmt{
				StmtList: ast.StmtList{
					{
						Body: &ast.CallStatement{
							Designator: asttest.NewDesignator(
								asttest.NewIdent("Write"),
							),
							ExprList: ast.ExprList{
								asttest.NewExpression(asttest.NewString("'Enter a value (0..9): '")),
							},
						},
					},
					{
						Body: &ast.CallStatement{
							Designator: asttest.NewDesignator(
								asttest.NewIdent("Readln"),
							),
							ExprList: ast.ExprList{
								asttest.NewExpression(asttest.NewQualId("I")),
							},
						},
					},
				},
				Condition: asttest.NewExpression(
					&ast.Term{
						Factor: &ast.Parentheses{
							Expression: &ast.Expression{
								SimpleExpression: asttest.NewSimpleExpression(asttest.NewQualId("I")),
								RelOpSimpleExpressions: ast.RelOpSimpleExpressions{
									{
										RelOp:            ">=",
										SimpleExpression: asttest.NewSimpleExpression(asttest.NewNumber("0")),
									},
								},
							},
						},
						MulOpFactors: []*ast.MulOpFactor{
							{
								MulOp: "AND",
								Factor: &ast.Parentheses{
									Expression: &ast.Expression{
										SimpleExpression: asttest.NewSimpleExpression(asttest.NewQualId("I")),
										RelOpSimpleExpressions: ast.RelOpSimpleExpressions{
											{
												RelOp:            "<=",
												SimpleExpression: asttest.NewSimpleExpression(asttest.NewNumber("9")),
											},
										},
									},
								},
							},
						},
					},
				),
			},
		},
	)

}
