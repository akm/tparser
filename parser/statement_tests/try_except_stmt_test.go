package statementtest

import (
	"testing"

	"github.com/akm/tparser/ast"
	"github.com/akm/tparser/ast/astcore"
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
							Body: &ast.CallStatement{
								Designator: asttest.NewDesignator("HandleException"),
							},
						},
					},
				},
			},
		},
	)

	stateToProgram := func(programName string, declSections ast.DeclSections, state *ast.Statement) *ast.Program {
		return &ast.Program{
			Ident: asttest.NewIdent(programName),
			ProgramBlock: &ast.ProgramBlock{
				Block: func() *ast.Block {
					return &ast.Block{
						DeclSections: declSections,
						CompoundStmt: &ast.CompoundStmt{
							StmtList: ast.StmtList{
								state,
							},
						},
					}
				}(),
			},
		}
	}

	typeDeclEZeroDivide := &ast.TypeDecl{
		Ident: asttest.NewIdent("EZeroDivide"),
		Type:  asttest.NewTypeId("Exception"),
	}
	typeDeclEOverflow := &ast.TypeDecl{
		Ident: asttest.NewIdent("EOverflow"),
		Type:  asttest.NewTypeId("Exception"),
	}
	typeDeclEMathError := &ast.TypeDecl{
		Ident: asttest.NewIdent("EMathError"),
		Type:  asttest.NewTypeId("Exception"),
	}

	runProgram(t,
		"with 1 ExceptionBlockHandler", true,
		[]rune(`PROGRAM OneExceptionBlockHandler;
type
	EZeroDivide = Exception;
begin
	try
		X := Y/Z;
	except
		on EZeroDivide do HandleZeroDivide;
	end;
end.
`),
		stateToProgram("OneExceptionBlockHandler",
			ast.DeclSections{
				ast.TypeSection{typeDeclEZeroDivide},
			},
			&ast.Statement{
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
								Type: asttest.NewTypeId(
									"EZeroDivide",
									astcore.NewDeclaration(typeDeclEZeroDivide.Ident, typeDeclEZeroDivide),
								),
								Statement: &ast.Statement{
									Body: &ast.CallStatement{
										Designator: asttest.NewDesignator(
											asttest.NewQualId("HandleZeroDivide"),
										),
									},
								},
							},
						},
					},
				},
			},
		),
	)

	runProgram(t,
		"with 3 ExceptionBlockHandlers without else", true,
		[]rune(`PROGRAM ThreeExceptionBlockHandlersWithoutElse;
type
	EZeroDivide = Exception;
	EOverflow = Exception;
	EMathError = Exception;
begin
	try
		DoSomething;
	except
		on EZeroDivide do HandleZeroDivide;
		on EOverflow do HandleOverflow;
		on EMathError do HandleMathError; 
	end;
end.
`),
		stateToProgram("ThreeExceptionBlockHandlersWithoutElse",
			ast.DeclSections{
				ast.TypeSection{
					typeDeclEZeroDivide,
					typeDeclEOverflow,
					typeDeclEMathError,
				},
			},
			&ast.Statement{
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
								Type: asttest.NewTypeId(
									"EZeroDivide",
									astcore.NewDeclaration(typeDeclEZeroDivide.Ident, typeDeclEZeroDivide),
								),
								Statement: &ast.Statement{
									Body: &ast.CallStatement{
										Designator: asttest.NewDesignator("HandleZeroDivide"),
									},
								},
							},
							{
								Type: asttest.NewTypeId(
									"EOverflow",
									astcore.NewDeclaration(typeDeclEOverflow.Ident, typeDeclEOverflow),
								),
								Statement: &ast.Statement{
									Body: &ast.CallStatement{
										Designator: asttest.NewDesignator("HandleOverflow"),
									},
								},
							},
							{
								Type: asttest.NewTypeId(
									"EMathError",
									astcore.NewDeclaration(typeDeclEMathError.Ident, typeDeclEMathError),
								),
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
		),
	)

	runProgram(t,
		"with 3 ExceptionBlockHandlers with else", true,
		[]rune(`PROGRAM ThreeExceptionBlockHandlersWithElse;
type
	EZeroDivide = Exception;
	EOverflow = Exception;
	EMathError = Exception;
begin
	try
		DoSomething;
	except
		on EZeroDivide do HandleZeroDivide;
		on EOverflow do HandleOverflow;
		on EMathError do HandleMathError;
	else
		HandleAllOthers;
	end;
end.
`),
		stateToProgram("ThreeExceptionBlockHandlersWithElse",
			ast.DeclSections{
				ast.TypeSection{
					typeDeclEZeroDivide,
					typeDeclEOverflow,
					typeDeclEMathError,
				},
			},
			&ast.Statement{
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
								Type: asttest.NewTypeId(
									"EZeroDivide",
									astcore.NewDeclaration(typeDeclEZeroDivide.Ident, typeDeclEZeroDivide),
								),
								Statement: &ast.Statement{
									Body: &ast.CallStatement{
										Designator: asttest.NewDesignator("HandleZeroDivide"),
									},
								},
							},
							{
								Type: asttest.NewTypeId(
									"EOverflow",
									astcore.NewDeclaration(typeDeclEOverflow.Ident, typeDeclEOverflow),
								),
								Statement: &ast.Statement{
									Body: &ast.CallStatement{
										Designator: asttest.NewDesignator("HandleOverflow"),
									},
								},
							},
							{
								Type: asttest.NewTypeId(
									"EMathError",
									astcore.NewDeclaration(typeDeclEMathError.Ident, typeDeclEMathError),
								),
								Statement: &ast.Statement{
									Body: &ast.CallStatement{
										Designator: asttest.NewDesignator("HandleMathError"),
									},
								},
							},
						},
						Else: ast.StmtList{
							{
								Body: &ast.CallStatement{
									Designator: asttest.NewDesignator("HandleAllOthers"),
								},
							},
						},
					},
				},
			},
		),
	)

	typeDeclException := &ast.TypeDecl{
		Ident: asttest.NewIdent("Exception"),
		Type:  asttest.NewTypeId("Error"),
	}

	runProgram(t,
		"with 1 ExceptionBlockHandler with ident",
		true,
		[]rune(`PROGRAM ExceptionBlockHandlerWithIdent;
type Exception = Error;
begin
	try
		DoSomething;
	except
		on E: Exception do ErrorDialog(E.Message1, E.HelpContext); // TODO Use Message instead of Message1
	end;
end.
`),
		stateToProgram("ExceptionBlockHandlerWithIdent",
			ast.DeclSections{
				ast.TypeSection{
					typeDeclException,
				},
			},
			&ast.Statement{
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
								Type: asttest.NewTypeId(
									"Exception",
									astcore.NewDeclaration(typeDeclException.Ident, typeDeclException),
								),
								Statement: &ast.Statement{
									Body: &ast.CallStatement{
										Designator: asttest.NewDesignator("ErrorDialog"),
										ExprList: ast.ExprList{
											asttest.NewExpression(
												&ast.Designator{
													QualId: asttest.NewQualId("E"), // TODO have astcore.Declarations to E
													Items: ast.DesignatorItems{
														asttest.NewDesignatorItemIdent("Message1"),
													},
												},
											),
											asttest.NewExpression(
												&ast.Designator{
													QualId: asttest.NewQualId("E"), // TODO have astcore.Declarations to E
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
		),
	)
}
