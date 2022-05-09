package statementtest

import (
	"testing"

	"github.com/akm/tparser/ast"
	"github.com/akm/tparser/ast/asttest"
	"github.com/akm/tparser/parser"
	"github.com/stretchr/testify/assert"
)

func TestWithStmt(t *testing.T) {
	run := func(name string, text []rune, expected *ast.Program) {
		t.Run(name, func(t *testing.T) {
			parser := parser.NewParser(&text)
			parser.NextToken()
			res, err := parser.ParseProgram()
			if assert.NoError(t, err) {
				asttest.ClearLocations(t, res)
				assert.Equal(t, expected, res)
				assert.Equal(t, expected.ProgramBlock.Block, res.ProgramBlock.Block)
			}

		})
	}

	// TODO Add TDate record type definition and check references to the fields
	// type TDate = record
	// 	Day: Integer;
	// 	Month: Integer;
	// 	Year: Integer;
	// end;

	run(
		"example1",
		[]rune(`PROGRAM WithExample1;
var OrderDate: TDate;
begin
	with OrderDate do
		if Month = 12 then begin
			Month := 1;
			Year := Year + 1;
		end
		else
			Month := Month + 1;
end.
`),
		&ast.Program{
			Ident: asttest.NewIdent("WithExample1"),
			ProgramBlock: func() *ast.ProgramBlock {
				decl := &ast.VarDecl{
					IdentList: asttest.NewIdentList("OrderDate"),
					Type:      asttest.NewTypeId("TDate"),
				}
				ref := decl.ToDeclarations()[0]

				return &ast.ProgramBlock{
					Block: &ast.Block{
						DeclSections: ast.DeclSections{
							ast.VarSection{decl},
						},
						CompoundStmt: &ast.CompoundStmt{
							StmtList: ast.StmtList{
								&ast.Statement{
									Body: &ast.WithStmt{
										Objects: ast.QualIds{
											asttest.NewQualId("OrderDate", ref),
										},
										Statement: &ast.Statement{
											Body: &ast.IfStmt{
												Condition: &ast.Expression{
													SimpleExpression: asttest.NewSimpleExpression("Month"), // TODO Use QualId with reference to decl
													RelOpSimpleExpressions: []*ast.RelOpSimpleExpression{
														{RelOp: "=", SimpleExpression: asttest.NewSimpleExpression(asttest.NewNumber("12"))},
													},
												},
												Then: &ast.Statement{
													Body: &ast.CompoundStmt{
														StmtList: ast.StmtList{
															{
																Body: &ast.AssignStatement{
																	Designator: asttest.NewDesignator("Month"), // TODO Use QualId with reference to decl
																	Expression: asttest.NewExpression(asttest.NewNumber("1")),
																},
															},
															{
																Body: &ast.AssignStatement{
																	Designator: asttest.NewDesignator("Year"), // TODO Use QualId with reference to decl
																	Expression: asttest.NewExpression(
																		&ast.SimpleExpression{
																			Term: &ast.Term{Factor: asttest.NewDesignatorFactor(asttest.NewIdent("Year"))}, // TODO Use QualId with reference to decl
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
													Body: &ast.AssignStatement{
														Designator: asttest.NewDesignator("Month"), // TODO Use QualId with reference to decl
														Expression: asttest.NewExpression(
															&ast.SimpleExpression{
																Term: &ast.Term{Factor: asttest.NewDesignatorFactor(asttest.NewIdent("Month"))}, // TODO Use QualId with reference to decl
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
								},
							},
						},
					},
				}
			}(),
		},
	)

}
