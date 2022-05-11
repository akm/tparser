package statementtest

import (
	"testing"

	"github.com/akm/tparser/ast"
	"github.com/akm/tparser/ast/asttest"
)

func TestTryExceptStmt(t *testing.T) {
	run := func(name string, text []rune, expected *ast.Statement) {
		runSatement(t, name, true, text, expected)
	}

	run(
		"without ExceptionBlockHandlers",
		[]rune(`
try
	DoSomething;
except
	HandleException;
end;
`),
		&ast.Statement{
			Body: &ast.RepeatStmt{
				StmtList: ast.StmtList{
					{
						Body: &ast.TryExceptStmt{
							Statements: ast.StmtList{
								{
									Body: &ast.CallStatement{
										Designator: asttest.NewDesignator("DoSomething"),
									},
								},
							},
							ExceptionBlock: &ast.ExceptionBlock{
								Else: ast.StmtList{
									{
										Body: &ast.AssignStatement{
											Designator: asttest.NewDesignator("HandleException"),
										},
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
		"with 1 ExceptionBlockHandler",
		[]rune(`
try
	X := Y/Z;
except
	on EZeroDivide do HandleZeroDivide;
end;
`),
		&ast.Statement{
			Body: &ast.RepeatStmt{
				StmtList: ast.StmtList{
					{
						Body: &ast.TryExceptStmt{
							Statements: ast.StmtList{
								{
									Body: &ast.AssignStatement{
										Designator: asttest.NewDesignator("X"),
										Expression: asttest.NewExpression(
											&ast.Term{
												Factor: asttest.NewDesignatorFactor(asttest.NewIdent("Y")),
												MulOpFactors: []*ast.MulOpFactor{
													{MulOp: "/", Factor: asttest.NewDesignatorFactor(asttest.NewIdent("Z"))},
												},
											},
										),
									},
								},
							},
							ExceptionBlock: &ast.ExceptionBlock{
								Handlers: ast.ExceptionBlockHandlers{
									{
										Type: asttest.NewTypeId("EZeroDivide"),
										Statement: &ast.Statement{
											Body: &ast.CallStatement{
												Designator: asttest.NewDesignator("HandleZeroDivide"),
											},
										},
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
		"with 3 ExceptionBlockHandlers without else",
		[]rune(`
try
	DoSomething;
except
	on EZeroDivide do HandleZeroDivide;
	on EOverflow do HandleOverflow;
	on EMathError do HandleMathError; 
end;
`),
		&ast.Statement{
			Body: &ast.RepeatStmt{
				StmtList: ast.StmtList{
					{
						Body: &ast.TryExceptStmt{
							Statements: ast.StmtList{
								{
									Body: &ast.CallStatement{
										Designator: asttest.NewDesignator("DoSomething"),
									},
								},
							},
							ExceptionBlock: &ast.ExceptionBlock{
								Handlers: ast.ExceptionBlockHandlers{
									{
										Type: asttest.NewTypeId("EZeroDivide"),
										Statement: &ast.Statement{
											Body: &ast.CallStatement{
												Designator: asttest.NewDesignator("HandleZeroDivide"),
											},
										},
									},
									{
										Type: asttest.NewTypeId("EOverflow"),
										Statement: &ast.Statement{
											Body: &ast.CallStatement{
												Designator: asttest.NewDesignator("HandleOverflow"),
											},
										},
									},
									{
										Type: asttest.NewTypeId("EMathError"),
										Statement: &ast.Statement{
											Body: &ast.CallStatement{
												Designator: asttest.NewDesignator("HandleMathError"),
											},
										},
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
		"with 3 ExceptionBlockHandlers with else",
		[]rune(`
try
	DoSomething;
except
	on EZeroDivide do HandleZeroDivide;
	on EOverflow do HandleOverflow;
	on EMathError do HandleMathError;
else
	HandleAllOthers;	
end;
`),
		&ast.Statement{
			Body: &ast.RepeatStmt{
				StmtList: ast.StmtList{
					{
						Body: &ast.TryExceptStmt{
							Statements: ast.StmtList{
								{
									Body: &ast.CallStatement{
										Designator: asttest.NewDesignator("DoSomething"),
									},
								},
							},
							ExceptionBlock: &ast.ExceptionBlock{
								Handlers: ast.ExceptionBlockHandlers{
									{
										Type: asttest.NewTypeId("EZeroDivide"),
										Statement: &ast.Statement{
											Body: &ast.CallStatement{
												Designator: asttest.NewDesignator("HandleZeroDivide"),
											},
										},
									},
									{
										Type: asttest.NewTypeId("EOverflow"),
										Statement: &ast.Statement{
											Body: &ast.CallStatement{
												Designator: asttest.NewDesignator("HandleOverflow"),
											},
										},
									},
									{
										Type: asttest.NewTypeId("EMathError"),
										Statement: &ast.Statement{
											Body: &ast.CallStatement{
												Designator: asttest.NewDesignator("HandleMathError"),
											},
										},
									},
								},
								Else: ast.StmtList{
									{
										Body: &ast.AssignStatement{
											Designator: asttest.NewDesignator("HandleAllOthers"),
										},
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
		"with 1 ExceptionBlockHandler with ident",
		[]rune(`
try
	DoSomething;
except
	on E: Exception do ErrorDialog(E.Message, E.HelpContext);
end;
`),
		&ast.Statement{
			Body: &ast.RepeatStmt{
				StmtList: ast.StmtList{
					{
						Body: &ast.TryExceptStmt{
							Statements: ast.StmtList{
								{
									Body: &ast.CallStatement{
										Designator: asttest.NewDesignator("DoSomething"),
									},
								},
							},
							ExceptionBlock: &ast.ExceptionBlock{
								Handlers: ast.ExceptionBlockHandlers{
									{
										Ident: asttest.NewIdent("E"),
										Type:  asttest.NewTypeId("Exception"),
										Statement: &ast.Statement{
											Body: &ast.CallStatement{
												Designator: asttest.NewDesignator("ErrorDialog"),
												ExprList: ast.ExprList{
													asttest.NewExpression(
														&ast.Designator{
															QualId: asttest.NewQualId("E"),
															Items: ast.DesignatorItems{
																asttest.NewDesignatorItemIdent("Message"),
															},
														},
													),
													asttest.NewExpression(
														&ast.Designator{
															QualId: asttest.NewQualId("E"),
															Items: ast.DesignatorItems{
																asttest.NewDesignatorItemIdent("HelpContext"),
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
					},
				},
			},
		},
	)
}
