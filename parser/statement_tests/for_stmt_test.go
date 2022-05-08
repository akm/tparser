package statementtest

import (
	"testing"

	"github.com/akm/tparser/ast"
	"github.com/akm/tparser/ast/asttest"
	"github.com/akm/tparser/parser"
	"github.com/stretchr/testify/assert"
)

func TestForStmt(t *testing.T) {
	run := func(name string, text []rune, expected *ast.Statement) {
		t.Run(name, func(t *testing.T) {
			parser := parser.NewParser(&text)
			parser.NextToken()
			res, err := parser.ParseStatement()
			if assert.NoError(t, err) {
				asttest.ClearLocations(t, res)
				if !assert.Equal(t, expected, res) {
					if !assert.Equal(t, expected.Body, res.Body) {
						expectFor := expected.Body.(*ast.ForStmt)
						actualFor := res.Body.(*ast.ForStmt)
						if !assert.Equal(t, expectFor.QualId, actualFor.QualId) {
						}
						if !assert.Equal(t, expectFor.Initial, actualFor.Initial) {
						}
						if !assert.Equal(t, expectFor.Terminal, actualFor.Terminal) {
						}
						if !assert.Equal(t, expectFor.Statement, actualFor.Statement) {
						}
					}
				}
			}
		})
	}

	run(
		"one line",
		[]rune(`for C := Red to Blue do Check(C);`),
		&ast.Statement{
			Body: &ast.ForStmt{
				QualId:   asttest.NewQualId("C"),
				Initial:  asttest.NewExpression(asttest.NewQualId("Red")),
				Terminal: asttest.NewExpression(asttest.NewQualId("Blue")),
				Statement: &ast.Statement{
					Body: &ast.CallStatement{
						Designator: asttest.NewDesignator("Check"),
						ExprList: ast.ExprList{
							asttest.NewExpression(asttest.NewQualId("C")),
						},
					},
				},
			},
		},
	)

	run(
		"example1",
		[]rune(`
for I := 2 to 63 do
	if Data[I] > Max then
		Max := Data[I];
`),
		&ast.Statement{
			Body: &ast.ForStmt{
				QualId:   asttest.NewQualId("I"),
				Initial:  asttest.NewExpression(asttest.NewNumber("2")),
				Terminal: asttest.NewExpression(asttest.NewNumber("63")),
				Statement: &ast.Statement{
					Body: &ast.IfStmt{
						Condition: &ast.Expression{
							SimpleExpression: asttest.NewSimpleExpression(
								&ast.Designator{
									QualId: asttest.NewQualId("Data"),
									Items: ast.DesignatorItems{
										ast.DesignatorItemExprList{asttest.NewExpression(asttest.NewQualId("I"))},
									},
								},
							),
							RelOpSimpleExpressions: []*ast.RelOpSimpleExpression{
								{
									RelOp:            ">",
									SimpleExpression: asttest.NewSimpleExpression(asttest.NewQualId("Max")),
								},
							},
						},
						Then: &ast.Statement{
							Body: &ast.AssignStatement{
								Designator: asttest.NewDesignator("Max"),
								Expression: asttest.NewExpression(
									&ast.Designator{
										QualId: asttest.NewQualId("Data"),
										Items: ast.DesignatorItems{
											ast.DesignatorItemExprList{
												asttest.NewExpression(asttest.NewQualId("I")),
											},
										},
									},
								),
							},
						},
					},
				},
			},
		},
	)

	run(
		"example2",
		[]rune(`
for I := ListBox1.Items.Count - 1 downto 0 do
	ListBox1.Items[I] := UpperCase(ListBox1.Items[I]);
`),
		&ast.Statement{
			Body: &ast.ForStmt{
				QualId: asttest.NewQualId("I"),
				Initial: asttest.NewExpression(
					&ast.SimpleExpression{
						Term: &ast.Term{Factor: asttest.NewDesignatorFactor(
							&ast.Designator{
								QualId: asttest.NewQualId("ListBox1"),
								Items: ast.DesignatorItems{
									asttest.NewDesignatorItemIdent(asttest.NewIdent("Items")),
									asttest.NewDesignatorItemIdent(asttest.NewIdent("Count")),
								},
							},
						)},
						AddOpTerms: []*ast.AddOpTerm{
							{AddOp: "-", Term: asttest.NewTerm(asttest.NewNumber("1"))},
						},
					},
				),
				Down:     true,
				Terminal: asttest.NewExpression(asttest.NewNumber("0")),
				Statement: &ast.Statement{
					Body: &ast.AssignStatement{
						Designator: &ast.Designator{
							QualId: asttest.NewQualId("ListBox1"),
							Items: ast.DesignatorItems{
								asttest.NewDesignatorItemIdent(asttest.NewIdent("Items")),
								ast.DesignatorItemExprList{
									asttest.NewExpression(asttest.NewQualId("I")),
								},
							},
						},
						Expression: asttest.NewExpression(
							&ast.DesignatorFactor{
								Designator: asttest.NewDesignator("UpperCase"),
								ExprList: ast.ExprList{
									asttest.NewExpression(
										&ast.Designator{
											QualId: asttest.NewQualId("ListBox1"),
											Items: ast.DesignatorItems{
												asttest.NewDesignatorItemIdent(asttest.NewIdent("Items")),
												ast.DesignatorItemExprList{
													asttest.NewExpression(asttest.NewQualId("I")),
												},
											},
										},
									),
								},
							},
						),
					},
				},
			},
		},
	)

	run(
		"example3",
		[]rune(`
for I := 1 to 10 do
	for J := 1 to 10 do
	begin
		X := 0;
		for K := 1 to 10 do
			X := X + Mat1[I, K] * Mat2[K, J];
		Mat[I, J] := X;
	end;
`),
		&ast.Statement{
			Body: &ast.ForStmt{
				QualId:   asttest.NewQualId("I"),
				Initial:  asttest.NewExpression(asttest.NewNumber("1")),
				Terminal: asttest.NewExpression(asttest.NewNumber("10")),
				Statement: &ast.Statement{
					Body: &ast.ForStmt{
						QualId:   asttest.NewQualId("J"),
						Initial:  asttest.NewExpression(asttest.NewNumber("1")),
						Terminal: asttest.NewExpression(asttest.NewNumber("10")),
						Statement: &ast.Statement{
							Body: &ast.CompoundStmt{
								StmtList: ast.StmtList{
									{
										Body: &ast.AssignStatement{
											Designator: asttest.NewDesignator("X"),
											Expression: asttest.NewExpression(asttest.NewNumber("0")),
										},
									},
									{
										Body: &ast.ForStmt{
											QualId:   asttest.NewQualId("K"),
											Initial:  asttest.NewExpression(asttest.NewNumber("1")),
											Terminal: asttest.NewExpression(asttest.NewNumber("10")),
											Statement: &ast.Statement{
												Body: &ast.AssignStatement{
													Designator: asttest.NewDesignator("X"),
													Expression: asttest.NewExpression(
														&ast.SimpleExpression{
															Term: &ast.Term{Factor: asttest.NewDesignatorFactor(asttest.NewIdent("X"))},
															AddOpTerms: []*ast.AddOpTerm{
																{
																	AddOp: "+",
																	Term: &ast.Term{
																		Factor: &ast.DesignatorFactor{
																			Designator: &ast.Designator{
																				QualId: asttest.NewQualId("Mat1"),
																				Items: ast.DesignatorItems{
																					ast.DesignatorItemExprList{
																						asttest.NewExpression(asttest.NewQualId("I")),
																						asttest.NewExpression(asttest.NewQualId("K")),
																					},
																				},
																			},
																		},
																		MulOpFactors: ast.MulOpFactors{
																			{
																				MulOp: "*",
																				Factor: &ast.DesignatorFactor{
																					Designator: &ast.Designator{
																						QualId: asttest.NewQualId("Mat2"),
																						Items: ast.DesignatorItems{
																							ast.DesignatorItemExprList{
																								asttest.NewExpression(asttest.NewQualId("K")),
																								asttest.NewExpression(asttest.NewQualId("J")),
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
													),
												},
											},
										},
									},
									{
										Body: &ast.AssignStatement{
											Designator: &ast.Designator{
												QualId: asttest.NewQualId("Mat"),
												Items: ast.DesignatorItems{
													ast.DesignatorItemExprList{
														asttest.NewExpression(asttest.NewQualId("I")),
														asttest.NewExpression(asttest.NewQualId("J")),
													},
												},
											},
											Expression: asttest.NewExpression(asttest.NewQualId("X")),
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
