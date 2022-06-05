package statementtest

import (
	"testing"

	"github.com/akm/tparser/ast"
	"github.com/akm/tparser/ast/asttest"
	"github.com/akm/tparser/parser/parsertest"
)

func TestIfStmt(t *testing.T) {
	parsertest.RunProgramTest(t,
		"Simple if statement",
		[]rune(`PROGRAM IfStmtSimple;
begin
	if Flag then
		writeln('OK');
end.
`),
		&ast.Program{
			Ident: asttest.NewIdent("IfStmtSimple"),
			ProgramBlock: &ast.ProgramBlock{
				Block: func() *ast.Block {

					return &ast.Block{
						Body: &ast.CompoundStmt{
							StmtList: ast.StmtList{
								&ast.Statement{
									Body: &ast.IfStmt{
										Condition: asttest.NewExpression(asttest.NewQualId("Flag")),
										Then: &ast.Statement{
											Body: &ast.CallStatement{
												Designator: asttest.NewDesignator(asttest.NewIdent("writeln")),
												ExprList: ast.ExprList{
													asttest.NewExpression(asttest.NewString("'OK'")),
												},
											},
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

	parsertest.RunProgramTest(t,
		"without ELSE",
		[]rune(`PROGRAM WihtoutElse;
var I,J: integer;
begin
	if J <> 0 then Result := I/J;
end.
`),
		&ast.Program{
			Ident: asttest.NewIdent("WihtoutElse"),
			ProgramBlock: &ast.ProgramBlock{
				Block: func() *ast.Block {
					decl := &ast.VarDecl{
						IdentList: asttest.NewIdentList([]string{"I", "J"}),
						Type:      asttest.NewOrdIdent("integer"),
					}
					declarations := decl.ToDeclarations()
					iDeclation := declarations[0]
					jDeclation := declarations[1]

					return &ast.Block{
						DeclSections: ast.DeclSections{ast.VarSection{decl}},
						Body: &ast.CompoundStmt{
							StmtList: ast.StmtList{
								&ast.Statement{
									Body: &ast.IfStmt{
										Condition: &ast.Expression{
											SimpleExpression: asttest.NewSimpleExpression(asttest.NewQualId("J", jDeclation)),
											RelOpSimpleExpressions: []*ast.RelOpSimpleExpression{
												{RelOp: "<>", SimpleExpression: asttest.NewSimpleExpression(asttest.NewNumber("0"))},
											},
										},
										Then: &ast.Statement{
											Body: &ast.AssignStatement{
												Designator: asttest.NewDesignator("Result"),
												Expression: asttest.NewExpression(
													&ast.Term{
														Factor: asttest.NewDesignatorFactor(asttest.NewQualId("I", iDeclation)),
														MulOpFactors: []*ast.MulOpFactor{
															{MulOp: "/", Factor: asttest.NewDesignatorFactor(asttest.NewQualId("J", jDeclation))},
														},
													},
												),
											},
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

	parsertest.RunProgramTest(t,
		"with ELSE",
		[]rune(`PROGRAM WihtElse;
var I,J: integer;
begin
	if J = 0 then
		Exit
	else
		Result := I/J;
end.
`),
		&ast.Program{
			Ident: asttest.NewIdent("WihtElse"),
			ProgramBlock: &ast.ProgramBlock{
				Block: func() *ast.Block {
					decl := &ast.VarDecl{
						IdentList: asttest.NewIdentList([]string{"I", "J"}),
						Type:      asttest.NewOrdIdent("integer"),
					}
					declarations := decl.ToDeclarations()
					iDeclation := declarations[0]
					jDeclation := declarations[1]

					return &ast.Block{
						DeclSections: ast.DeclSections{ast.VarSection{decl}},
						Body: &ast.CompoundStmt{
							StmtList: ast.StmtList{
								&ast.Statement{
									Body: &ast.IfStmt{
										Condition: &ast.Expression{
											SimpleExpression: asttest.NewSimpleExpression(asttest.NewQualId("J", jDeclation)),
											RelOpSimpleExpressions: []*ast.RelOpSimpleExpression{
												{RelOp: "=", SimpleExpression: asttest.NewSimpleExpression(asttest.NewNumber("0"))},
											},
										},
										Then: &ast.Statement{
											Body: &ast.CallStatement{
												Designator: asttest.NewDesignator("Exit"),
											},
										},
										Else: &ast.Statement{
											Body: &ast.AssignStatement{
												Designator: asttest.NewDesignator("Result"),
												Expression: asttest.NewExpression(
													&ast.Term{
														Factor: asttest.NewDesignatorFactor(asttest.NewQualId("I", iDeclation)),
														MulOpFactors: []*ast.MulOpFactor{
															{MulOp: "/", Factor: asttest.NewDesignatorFactor(asttest.NewQualId("J", jDeclation))},
														},
													},
												),
											},
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

	parsertest.RunProgramTest(t,
		"with ELSE IF",
		[]rune(`PROGRAM WihtElseIf;
begin
	if J = 0 then
	begin
		Result := I/J;
		Count := Count + 1;
	end
	else if Count = Last then
		Done := True
	else
		Exit;
end.
`),
		&ast.Program{
			Ident: asttest.NewIdent("WihtElseIf"),
			ProgramBlock: &ast.ProgramBlock{
				Block: func() *ast.Block {
					return &ast.Block{
						Body: &ast.CompoundStmt{
							StmtList: ast.StmtList{
								&ast.Statement{
									Body: &ast.IfStmt{
										Condition: &ast.Expression{
											SimpleExpression: asttest.NewSimpleExpression("J"),
											RelOpSimpleExpressions: []*ast.RelOpSimpleExpression{
												{RelOp: "=", SimpleExpression: asttest.NewSimpleExpression(asttest.NewNumber("0"))},
											},
										},
										Then: &ast.Statement{
											Body: &ast.CompoundStmt{
												StmtList: ast.StmtList{
													{
														Body: &ast.AssignStatement{
															Designator: asttest.NewDesignator("Result"),
															Expression: asttest.NewExpression(
																&ast.Term{
																	Factor: asttest.NewDesignatorFactor("I"),
																	MulOpFactors: []*ast.MulOpFactor{
																		{MulOp: "/", Factor: asttest.NewDesignatorFactor("J")},
																	},
																},
															),
														},
													},
													{
														Body: &ast.AssignStatement{
															Designator: asttest.NewDesignator("Count"),
															Expression: asttest.NewExpression(
																&ast.SimpleExpression{
																	Term: &ast.Term{Factor: asttest.NewDesignatorFactor(asttest.NewIdent("Count"))},
																	AddOpTerms: []*ast.AddOpTerm{
																		{AddOp: "+", Term: asttest.NewTerm(asttest.NewNumber("1"))},
																	},
																},
															),
														},
													},
												},
											},
										},
										Else: &ast.Statement{
											Body: &ast.IfStmt{
												Condition: &ast.Expression{
													SimpleExpression: asttest.NewSimpleExpression("Count"),
													RelOpSimpleExpressions: []*ast.RelOpSimpleExpression{
														{
															RelOp:            "=",
															SimpleExpression: asttest.NewSimpleExpression("Last"),
														},
													},
												},
												Then: &ast.Statement{
													Body: &ast.AssignStatement{
														Designator: asttest.NewDesignator("Done"),
														Expression: asttest.NewExpression(ast.NewValueFactor("True")),
													},
												},
												Else: &ast.Statement{
													Body: &ast.CallStatement{
														Designator: asttest.NewDesignator("Exit"),
													},
												},
											},
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
