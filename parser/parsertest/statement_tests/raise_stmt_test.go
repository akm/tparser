package statementtest

import (
	"testing"

	"github.com/akm/tparser/ast"
	"github.com/akm/tparser/ast/astcore"
	"github.com/akm/tparser/ast/asttest"
	"github.com/akm/tparser/parser/parsertest"
	"github.com/stretchr/testify/assert"
)

func TestRaiseStmt(t *testing.T) {
	runWithFunc := func(t *testing.T, name string, text []rune, expected *ast.FunctionDecl) {
		t.Run(name, func(t *testing.T) {
			parser := parsertest.NewTestParser(&text)
			parser.NextToken()
			res, err := parser.ParseProcedureDeclSection()
			if assert.NoError(t, err) {
				asttest.ClearLocations(t, res)
				assert.Equal(t, expected, res)
			}
		})
	}

	declPath := &ast.FormalParm{
		Opt: &ast.FpoConst,
		Parameter: &ast.Parameter{
			IdentList: asttest.NewIdentList("Path"),
			Type: &ast.ParameterType{
				Type: asttest.NewStringType("string"),
			},
		},
	}
	declI := &ast.VarDecl{
		IdentList: asttest.NewIdentList("I"),
		Type:      asttest.NewOrdIdent("Integer"),
	}
	declSearchRec := &ast.VarDecl{
		IdentList: asttest.NewIdentList("SearchRec"),
		Type:      asttest.NewTypeId("TSearchRec"),
	}

	runWithFunc(t,
		"basic",
		[]rune(`
function GetFileList(const Path: string): TStringList;
var
	I: Integer;
	SearchRec: TSearchRec;
begin
	Result := TStringList.Create;
	try
		I := FindFirst(Path, 0, SearchRec);
		while I = 0 do
		begin
			Result.Add(SearchRec.Name);
			I := FindNext(SearchRec);
		end;
	except
		Result.Free;
		raise;
	end;
end;
`),
		&ast.FunctionDecl{
			FunctionHeading: &ast.FunctionHeading{
				Type:             ast.FtFunction,
				Ident:            asttest.NewIdent("GetFileList"),
				FormalParameters: ast.FormalParameters{declPath},
				ReturnType:       &ast.TypeId{Ident: asttest.NewIdent("TStringList")},
			},
			Block: &ast.Block{
				DeclSections: ast.DeclSections{
					ast.VarSection{declI, declSearchRec},
				},
				Body: &ast.CompoundStmt{
					StmtList: ast.StmtList{
						{
							// Result := TStringList.Create;
							Body: &ast.AssignStatement{
								Designator: asttest.NewDesignator("Result"),
								Expression: asttest.NewExpression(
									&ast.DesignatorFactor{
										Designator: &ast.Designator{
											QualId: asttest.NewQualId("TStringList"),
											Items: ast.DesignatorItems{
												asttest.NewDesignatorItemIdent("Create"),
											},
										},
									},
								),
							},
						},
						{
							Body: &ast.TryExceptStmt{
								Statements: ast.StmtList{
									{
										// I := FindFirst(Path, 0, SearchRec);
										Body: &ast.AssignStatement{
											Designator: asttest.NewDesignator(
												asttest.NewQualId(
													"I",
													astcore.NewDeclaration(declI.IdentList[0], declI),
												),
											),
											Expression: asttest.NewExpression(
												&ast.DesignatorFactor{
													Designator: &ast.Designator{
														QualId: asttest.NewQualId("FindFirst"),
													},
													ExprList: ast.ExprList{
														asttest.NewExpression(
															asttest.NewQualId(
																"Path",
																astcore.NewDeclaration(declPath.IdentList[0], declPath),
															),
														),
														asttest.NewExpression(asttest.NewNumber("0")),
														asttest.NewExpression(
															asttest.NewQualId(
																"SearchRec",
																astcore.NewDeclaration(declSearchRec.IdentList[0], declSearchRec),
															),
														),
													},
												},
											),
										},
									},
									{
										Body: &ast.WhileStmt{
											Condition: &ast.Expression{
												SimpleExpression: ast.NewSimpleExpression(
													asttest.NewQualId(
														"I",
														astcore.NewDeclaration(declI.IdentList[0], declI),
													),
												),
												RelOpSimpleExpressions: []*ast.RelOpSimpleExpression{
													{RelOp: "=", SimpleExpression: asttest.NewSimpleExpression(asttest.NewNumber("0"))},
												},
											},
											Statement: &ast.Statement{
												Body: &ast.CompoundStmt{
													StmtList: ast.StmtList{
														// Result.Add(SearchRec.Name);
														{
															Body: &ast.CallStatement{
																Designator: &ast.Designator{
																	QualId: asttest.NewQualId("Result"),
																	Items: ast.DesignatorItems{
																		asttest.NewDesignatorItemIdent("Add"),
																	},
																},
																ExprList: ast.ExprList{
																	asttest.NewExpression(
																		&ast.Designator{
																			QualId: asttest.NewQualId(
																				"SearchRec",
																				astcore.NewDeclaration(declSearchRec.IdentList[0], declSearchRec),
																			),
																			Items: ast.DesignatorItems{
																				asttest.NewDesignatorItemIdent("Name"),
																			},
																		},
																	),
																},
															},
														},
														{
															// I := FindNext(SearchRec);
															Body: &ast.AssignStatement{
																Designator: asttest.NewDesignator(
																	asttest.NewQualId(
																		"I",
																		astcore.NewDeclaration(declI.IdentList[0], declI),
																	),
																),
																Expression: asttest.NewExpression(
																	&ast.DesignatorFactor{
																		Designator: &ast.Designator{
																			QualId: asttest.NewQualId("FindNext"),
																		},
																		ExprList: ast.ExprList{
																			asttest.NewExpression(
																				asttest.NewQualId(
																					"SearchRec",
																					astcore.NewDeclaration(declSearchRec.IdentList[0], declSearchRec),
																				),
																			),
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
								ExceptionBlock: &ast.ExceptionBlock{
									Else: ast.StmtList{
										{
											// Result.Free;
											Body: &ast.CallStatement{
												Designator: &ast.Designator{
													QualId: asttest.NewQualId("Result"),
													Items: ast.DesignatorItems{
														asttest.NewDesignatorItemIdent("Free"),
													},
												},
											},
										},
										{
											// raise;
											Body: &ast.RaiseStmt{},
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

	// See https://docwiki.embarcadero.com/RADStudio/Sydney/en/Exceptions_(Delphi)#Raising_and_Handling_Exceptions
	// raise Exception.Create('Missing parameter') at @MyFunction;
	parsertest.RunStatementTest(t, "raise with at",
		[]rune(
			`raise Exception.Create('Missing parameter') at @MyFunction;`,
		),
		&ast.Statement{
			Body: &ast.RaiseStmt{
				Object: asttest.NewExpression(
					asttest.NewExpression(
						&ast.DesignatorFactor{
							Designator: &ast.Designator{
								QualId: asttest.NewQualId("Exception"),
								Items: ast.DesignatorItems{
									asttest.NewDesignatorItemIdent("Create"),
								},
							},
							ExprList: ast.ExprList{
								asttest.NewExpression(asttest.NewString("'Missing parameter'")),
							},
						},
					),
				),
				Address: asttest.NewExpression(
					&ast.Address{
						Designator: asttest.NewDesignator(asttest.NewQualId("MyFunction")),
					},
				),
			},
		},
	)

	// TODO implement test after class support
	// type
	// 	ETrigError = class(EMathError);
	// function Tan(X: Extended): Extended;
	// begin
	// 	try
	// 		Result := Sin(X) / Cos(X);
	// 	except
	// 		on EMathError do
	// 			raise ETrigError.Create('Invalid argument to Tan');
	// 	end;
	// end;

}
