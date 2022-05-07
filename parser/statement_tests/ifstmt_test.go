package statementtest

import (
	"testing"

	"github.com/akm/tparser/ast"
	"github.com/akm/tparser/ast/asttest"
	"github.com/akm/tparser/parser"
	"github.com/stretchr/testify/assert"
)

func TestIfStmt(t *testing.T) {
	run := func(name string, clearLocations bool, text []rune, expected *ast.Program) {
		t.Run(name, func(t *testing.T) {
			parser := parser.NewParser(&text)
			parser.NextToken()
			res, err := parser.ParseProgram()
			if assert.NoError(t, err) {
				if clearLocations {
					asttest.ClearLocations(t, res)
				}
				if !assert.Equal(t, expected, res) {
					if !assert.Equal(t, expected.ProgramBlock.Block, res.ProgramBlock.Block) {
						if !assert.Equal(t, expected.ProgramBlock.Block.CompoundStmt, res.ProgramBlock.Block.CompoundStmt) {
							if !assert.Equal(t, expected.ProgramBlock.Block.CompoundStmt.StmtList, res.ProgramBlock.Block.CompoundStmt.StmtList) {
								if !assert.Equal(t, expected.ProgramBlock.Block.CompoundStmt.StmtList[0], res.ProgramBlock.Block.CompoundStmt.StmtList[0]) {
									expectBody := expected.ProgramBlock.Block.CompoundStmt.StmtList[0].Body.(*ast.IfStmt)
									actualBody := res.ProgramBlock.Block.CompoundStmt.StmtList[0].Body.(*ast.IfStmt)
									if !assert.Equal(t, expectBody, actualBody) {
										assert.Equal(t, expectBody.Condition, actualBody.Condition)
										if !assert.Equal(t, expectBody.Then, actualBody.Then) {
											expectBody := expectBody.Then.Body.(*ast.CompoundStmt)
											actualBody := actualBody.Then.Body.(*ast.CompoundStmt)
											if !assert.Equal(t, expectBody, actualBody) {
												if !assert.Equal(t, expectBody.StmtList, actualBody.StmtList) {
													if !assert.Equal(t, expectBody.StmtList[0], actualBody.StmtList[0]) {

													}
													if !assert.Equal(t, expectBody.StmtList[1], actualBody.StmtList[1]) {
														if !assert.Equal(t, expectBody.StmtList[1].Body, actualBody.StmtList[1].Body) {

														}
													}
												}
											}
										}
										assert.Equal(t, expectBody.Else, actualBody.Else)
									}
								}
							}
						}
					}
				}
			}
		})
	}

	run(
		"Simple if statement", true,
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
						CompoundStmt: &ast.CompoundStmt{
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

	run(
		"without ELSE", true,
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
						CompoundStmt: &ast.CompoundStmt{
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

	run(
		"with ELSE", true,
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
						CompoundStmt: &ast.CompoundStmt{
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

	run(
		"with ELSE IF", true,
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
						CompoundStmt: &ast.CompoundStmt{
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
